/**
 * Generates favicon.ico from favicon.svg for browsers that require ICO.
 * Run: node scripts/generate-favicon.mjs
 */
import fs from "fs";
import path from "path";
import { fileURLToPath } from "url";
import sharp from "sharp";
import toIco from "to-ico";

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const publicDir = path.join(__dirname, "..", "public");
const svgPath = path.join(publicDir, "favicon.svg");
const icoPath = path.join(publicDir, "favicon.ico");

const svgBuffer = fs.readFileSync(svgPath);
// Use explicit stroke so sharp renders correctly (currentColor may not resolve in Node)
const svgWithColor = Buffer.from(
  svgBuffer.toString().replace(/currentColor/g, "#333333")
);

async function main() {
  const sizes = [16, 32];
  const pngBuffers = await Promise.all(
    sizes.map((size) =>
      sharp(svgWithColor).resize(size, size).png().toBuffer()
    )
  );
  const icoBuffer = await toIco(pngBuffers);
  fs.writeFileSync(icoPath, icoBuffer);
  console.log("Wrote", icoPath);
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
