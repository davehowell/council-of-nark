// Verify the production build and capture every slide for human visual review.
import { createServer } from 'node:http';
import { readFileSync, existsSync, statSync, mkdirSync, mkdtempSync } from 'node:fs';
import { join, extname, resolve } from 'node:path';
import { tmpdir } from 'node:os';
import { chromium } from 'playwright-chromium';
import { decks } from './decks.mjs';

const root = new URL('../dist/', import.meta.url).pathname;
const output = process.argv[2] || mkdtempSync(join(tmpdir(), 'council-slide-review-'));
mkdirSync(output, { recursive: true });
const deckIndex = readFileSync(join(root, 'decks', 'part-1', 'index.html'), 'utf8');
const sitePrefix = deckIndex.match(/src="([^"]*)\/decks\/part-1\/assets\//)?.[1] || '';
const types = { '.html': 'text/html', '.js': 'application/javascript', '.css': 'text/css', '.png': 'image/png', '.jpg': 'image/jpeg', '.jpeg': 'image/jpeg', '.svg': 'image/svg+xml', '.pdf': 'application/pdf' };
const server = createServer((req, res) => {
  const requestedPath = decodeURIComponent(new URL(req.url, 'http://localhost').pathname);
  const pathname = sitePrefix && requestedPath.startsWith(`${sitePrefix}/`) ? requestedPath.slice(sitePrefix.length) : requestedPath;
  let file = resolve(root, `.${pathname}`);
  if (!file.startsWith(resolve(root) + '/')) file = join(root, 'index.html');
  if (existsSync(file) && statSync(file).isDirectory()) file = join(file, 'index.html');
  if (!existsSync(file)) {
    const slug = pathname.match(/^\/decks\/(part-[1-4])\/(?:\d+|presenter\/\d+)$/)?.[1];
    if (slug) file = join(root, 'decks', slug, 'index.html');
  }
  if (!existsSync(file)) { res.writeHead(404); res.end(); return; }
  res.setHeader('Content-Type', types[extname(file)] || 'text/plain');
  res.end(readFileSync(file));
});
await new Promise(resolve => server.listen(0, '127.0.0.1', resolve));
const origin = `http://127.0.0.1:${server.address().port}${sitePrefix}`;
const browser = await chromium.launch();
const errors = [];
try {
  const page = await browser.newPage({ viewport: { width: 1440, height: 900 }, deviceScaleFactor: 1 });
  page.on('pageerror', e => {
    // Headless Chromium may refuse Slidev's optional keep-screen-awake request.
    // This has no effect on slide rendering; all other page errors fail QA.
    if (e.message !== 'Wake Lock permission request denied') errors.push(e.message);
  });
  page.on('response', r => { if (r.status() >= 400) errors.push(`${r.status()}: ${r.url()}`); });
  for (const deck of decks) {
    for (let n = 1; n <= deck.slides; n++) {
      await page.goto(`${origin}/decks/${deck.slug}/${n}`, { waitUntil: 'networkidle' });
      await page.locator('.slidev-layout:visible').first().waitFor();
      await page.evaluate(() => document.fonts.ready);
      const slide = page.locator('.slidev-layout:visible').first();
      const problems = await slide.evaluate(el => {
        const outer = el.getBoundingClientRect();
        const footer = el.querySelector('.foot')?.getBoundingClientRect();
        return [...el.querySelectorAll('h1,h2,p,li,table,.pair')].flatMap(node => {
          const r = node.getBoundingClientRect();
          if (r.width === 0 || r.height === 0) return [];
          const outside = r.left < outer.left - 1 || r.right > outer.right + 1 || r.top < outer.top - 1 || r.bottom > outer.bottom + 1;
          const footerOverlap = footer && r.bottom > footer.top && r.top < footer.bottom;
          return outside || footerOverlap ? [node.textContent.trim().slice(0, 100)] : [];
        });
      });
      if (problems.length) errors.push(`${deck.slug}/${n} overflow: ${problems.join('; ')}`);
      await slide.screenshot({ path: join(output, `${deck.slug}-${n}.png`) });
    }
    await page.goto(`${origin}/downloads/${deck.slug}-notes.html`, { waitUntil: 'networkidle' });
    if (await page.locator('section').count() !== deck.slides) errors.push(`${deck.slug}: notes do not contain ${deck.slides} sections`);
  }
  const slideCount = decks.reduce((total, deck) => total + deck.slides, 0);
  console.log(`Captured ${slideCount} slides: ${output}`);
  if (errors.length) throw new Error([...new Set(errors)].join('\n'));
  console.log('All slides and notes loaded; no content overflow, footer overlap, page errors or failed requests.');
} finally {
  await browser.close();
  server.close();
}
