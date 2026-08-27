import { cpSync, mkdirSync, rmSync, writeFileSync } from "node:fs";
import { execFileSync } from "node:child_process";
import { dirname, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";

const presentationsRoot = resolve(dirname(fileURLToPath(import.meta.url)), "..");
const repositoryRoot = resolve(presentationsRoot, "..");
const outputRoot = join(presentationsRoot, "dist");
const siteRoot = join(presentationsRoot, "site");
const slidev = join(presentationsRoot, "node_modules", ".bin", "slidev");
const configuredBase = (process.env.SITE_BASE || "").trim();
const siteBase = configuredBase === "/" ? "" : configuredBase.replace(/^\/*|\/*$/g, "");
const basePrefix = siteBase ? `/${siteBase}` : "";

const decks = [
  { slug: "part-1", entry: "part-1/slides.md", pdf: "part-1/council-of-nark.pdf" },
  { slug: "part-2", entry: "part-2/slides.md", pdf: "part-2/put-the-council-on-trial.pdf" },
  { slug: "part-3", entry: "part-3/slides.md", pdf: "part-3/the-experiment-fought-back.pdf" },
];

rmSync(outputRoot, { recursive: true, force: true });
mkdirSync(outputRoot, { recursive: true });
cpSync(siteRoot, outputRoot, { recursive: true });
mkdirSync(join(outputRoot, "assets"), { recursive: true });
mkdirSync(join(outputRoot, "downloads"), { recursive: true });
cpSync(join(presentationsRoot, "public", "Nark-council.png"), join(outputRoot, "assets", "council.png"));

for (const deck of decks) {
  const output = join(outputRoot, "decks", deck.slug);
  const base = `${basePrefix}/decks/${deck.slug}/`;
  execFileSync(
    slidev,
    ["build", join(presentationsRoot, deck.entry), "--base", base, "--out", output],
    { cwd: presentationsRoot, stdio: "inherit" },
  );
  cpSync(join(presentationsRoot, deck.pdf), join(outputRoot, "downloads", deck.pdf.split("/").at(-1)));
}

writeFileSync(join(outputRoot, ".nojekyll"), "");
writeFileSync(
  join(outputRoot, "build.json"),
  `${JSON.stringify({ schema_version: 1, repository: "davehowell/council-of-nark", source: process.env.GITHUB_SHA || "local" }, null, 2)}\n`,
);

console.log(`Built Pages site at ${outputRoot}`);
console.log(`Base path: ${basePrefix || "/"}`);
console.log(`Source root: ${repositoryRoot}`);
