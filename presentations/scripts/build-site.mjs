import { createHash } from "node:crypto";
import { cpSync, mkdirSync, readFileSync, rmSync, writeFileSync } from "node:fs";
import { execFileSync } from "node:child_process";
import { dirname, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import { decks } from "./decks.mjs";
import { writeNotes } from "./speaker-notes.mjs";

const presentationsRoot = resolve(dirname(fileURLToPath(import.meta.url)), "..");
const repositoryRoot = resolve(presentationsRoot, "..");
const outputRoot = join(presentationsRoot, "dist");
const siteRoot = join(presentationsRoot, "site");
const slidev = join(presentationsRoot, "node_modules", ".bin", "slidev");
const configuredBase = (process.env.SITE_BASE || "").trim();
const siteBase = configuredBase === "/" ? "" : configuredBase.replace(/^\/*|\/*$/g, "");
const basePrefix = siteBase ? `/${siteBase}` : "";

rmSync(outputRoot, { recursive: true, force: true });
mkdirSync(outputRoot, { recursive: true });
cpSync(siteRoot, outputRoot, { recursive: true });
mkdirSync(join(outputRoot, "assets"), { recursive: true });
mkdirSync(join(outputRoot, "downloads"), { recursive: true });
cpSync(join(presentationsRoot, "public", "Nark-council.png"), join(outputRoot, "assets", "council.png"));

// GitHub Pages can briefly serve new HTML with a cached prior asset. Content-derived
// query strings keep each deployment's HTML, CSS, JavaScript, and hero image together.
const versionedAssets = ["site.css", "site.js", "council.png"];
const assetVersions = Object.fromEntries(versionedAssets.map((asset) => {
  const content = readFileSync(join(outputRoot, "assets", asset));
  return [asset, createHash("sha256").update(content).digest("hex").slice(0, 12)];
}));
for (const page of ["index.html", "404.html", "results/index.html", "timeline/index.html"]) {
  const path = join(outputRoot, page);
  let html = readFileSync(path, "utf8");
  for (const [asset, version] of Object.entries(assetVersions)) {
    html = html.replaceAll(asset, `${asset}?v=${version}`);
  }
  writeFileSync(path, html);
}

for (const deck of decks) {
  const output = join(outputRoot, "decks", deck.slug);
  const base = `${basePrefix}/decks/${deck.slug}/`;
  execFileSync(
    slidev,
    ["build", join(presentationsRoot, deck.slug, "slides.md"), "--base", base, "--out", output],
    { cwd: presentationsRoot, stdio: "inherit" },
  );
  // Export from the same source as the web deck; never publish a stale PDF.
  execFileSync(slidev, ["export", join(presentationsRoot, deck.slug, "slides.md"), "--output", join(outputRoot, "downloads", deck.pdf)], { cwd: presentationsRoot, stdio: "inherit" });
  for (const alias of deck.pdfAliases || []) {
    cpSync(join(outputRoot, "downloads", deck.pdf), join(outputRoot, "downloads", alias));
  }
}
writeNotes(join(outputRoot, "downloads"));

writeFileSync(join(outputRoot, ".nojekyll"), "");
writeFileSync(
  join(outputRoot, "build.json"),
  `${JSON.stringify({ schema_version: 1, repository: "davehowell/council-of-nark", source: process.env.GITHUB_SHA || "local" }, null, 2)}\n`,
);

console.log(`Built Pages site at ${outputRoot}`);
console.log(`Base path: ${basePrefix || "/"}`);
console.log(`Source root: ${repositoryRoot}`);
