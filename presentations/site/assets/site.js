const reduceMotion = window.matchMedia("(prefers-reduced-motion: reduce)").matches;
const canvas = document.createElement("canvas");
const grid = document.createElement("div");
const scanlines = document.createElement("div");

canvas.className = "starfield";
grid.className = "synth-grid";
scanlines.className = "scanlines";
canvas.setAttribute("aria-hidden", "true");
grid.setAttribute("aria-hidden", "true");
scanlines.setAttribute("aria-hidden", "true");
document.body.prepend(scanlines);
document.body.prepend(grid);
document.body.prepend(canvas);

const context = canvas.getContext("2d");
let stars = [];
let animationFrame;

function resizeStars() {
  const scale = Math.min(window.devicePixelRatio || 1, 2);
  canvas.width = Math.floor(window.innerWidth * scale);
  canvas.height = Math.floor(window.innerHeight * scale);
  canvas.style.width = `${window.innerWidth}px`;
  canvas.style.height = `${window.innerHeight}px`;
  context.setTransform(scale, 0, 0, scale, 0, 0);

  const count = Math.min(320, Math.floor((window.innerWidth * window.innerHeight) / 5200));
  stars = Array.from({ length: count }, () => ({
    x: Math.random() * window.innerWidth,
    y: Math.random() * window.innerHeight,
    radius: Math.random() * 1.1 + 0.25,
    alpha: Math.random() * 0.45 + 0.2,
    phase: Math.random() * Math.PI * 2,
    speed: Math.random() * 0.0008 + 0.00035,
  }));
}

function drawStars(time = 0) {
  context.clearRect(0, 0, window.innerWidth, window.innerHeight);
  for (const star of stars) {
    const flicker = reduceMotion ? 0.75 : 0.65 + Math.sin(time * star.speed + star.phase) * 0.35;
    context.fillStyle = `rgba(190, 225, 255, ${star.alpha * flicker})`;
    context.beginPath();
    context.arc(star.x, star.y, star.radius, 0, Math.PI * 2);
    context.fill();
  }

  if (!reduceMotion) animationFrame = window.requestAnimationFrame(drawStars);
}

function resetStars() {
  window.cancelAnimationFrame(animationFrame);
  resizeStars();
  drawStars();
}

window.addEventListener("resize", resetStars);
resizeStars();
drawStars();
