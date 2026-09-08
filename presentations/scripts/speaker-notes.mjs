import { readFileSync, writeFileSync, mkdirSync } from 'node:fs';
import { join, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';
import { decks } from './decks.mjs';

const root = fileURLToPath(new URL('..', import.meta.url));
const escape = value => value.replaceAll('&', '&amp;').replaceAll('<', '&lt;').replaceAll('>', '&gt;').replaceAll('"', '&quot;');

export function writeNotes(output) {
  mkdirSync(output, { recursive: true });
  const all = [];
  for (const deck of decks) {
    const source = readFileSync(join(root, deck.slug, 'slides.md'), 'utf8');
    const sections = source.replace(/^---\n[\s\S]*?\n---\n/, '').split(/\n---\n/);
    if (sections.length !== 6) throw new Error(`${deck.slug}: expected six slides`);
    const notes = sections.map((slide, i) => {
      const title = slide.match(/^# (.+)$/m)?.[1].replaceAll('<br>', ' ');
      const note = slide.match(/<!--\s*([\s\S]*?)\s*-->\s*$/)?.[1];
      if (!title || !note || !note.includes('[Sources]') || !note.includes('[Time:')) throw new Error(`${deck.slug}/${i + 1}: missing title, timing, sources or notes`);
      return { title, note, number: i + 1 };
    });
    const spoken = notes.map(n => n.note.split('[Sources]')[0].replace(/\[Time:[^\]]+\]/g, '')).join(' ');
    const words = spoken.trim().split(/\s+/).length;
    const markdown = `# ${deck.title} — speaker notes\n\nSix slides · approximately five minutes · ${words} spoken words.\n\nEdit the closing comment in each slide in ${deck.slug}/slides.md; these notes are generated from that source.\n\n` + notes.map(n => `## ${n.number}. ${n.title}\n\n${n.note}\n`).join('\n');
    writeFileSync(join(output, `${deck.slug}-notes.md`), markdown);
    const body = notes.map(n => `<section><h2>${n.number}. ${escape(n.title)}</h2>${n.note.split(/\n\n/).map(p => p.startsWith('[Sources]') ? `<details><summary>Sources</summary>${p.split('\n').slice(1).map(u => `<p><a href="${escape(u)}">${escape(u)}</a></p>`).join('')}</details>` : `<p>${escape(p)}</p>`).join('')}</section>`).join('');
    const html = `<!doctype html><html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>${escape(deck.title)} — speaker notes</title><style>body{font:20px/1.65 Arial,sans-serif;color:#172033;max-width:780px;margin:48px auto;padding:0 24px}h1{font-size:40px;line-height:1.15}h2{font-size:26px}section{border-top:1px solid #ccd3df;margin-top:40px;padding-top:20px}a{color:#315683;overflow-wrap:anywhere}details{font-size:14px}nav{font-size:16px}@media print{nav,details{display:none}section{break-inside:avoid}body{font-size:13pt;margin:0}}</style></head><body><nav><a href="../">All talks</a> · <a href="${deck.slug}-notes.md">Markdown notes</a> · <a href="${deck.pdf}">Slides PDF</a></nav><h1>${escape(deck.title)}</h1><p>Speaker notes · 6 slides · approximately 5 minutes · ${words} words</p>${body}</body></html>`;
    // Printable HTML navigation targets the published downloads directory.
    // Local authoring copies are Markdown, with slides.md as the authority.
    if (resolve(output) !== resolve(root, 'notes')) writeFileSync(join(output, `${deck.slug}-notes.html`), html);
    all.push(markdown);
    console.log(`${deck.slug}: 6 slides, ${words} spoken words`);
  }
  writeFileSync(join(output, 'speaker-notes.md'), all.join('\n---\n\n'));
}

if (process.argv[1] && resolve(process.argv[1]) === fileURLToPath(import.meta.url)) writeNotes(resolve(process.argv[2] || join(root, 'notes')));
