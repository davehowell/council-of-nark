const events = Array.isArray(window.COUNCIL_TIMELINE) ? window.COUNCIL_TIMELINE : [];
const timeline = document.querySelector("#timeline");
const filters = [...document.querySelectorAll("[data-filter]")];
const presentButton = document.querySelector("#present-timeline");
const dialog = document.querySelector("#story-dialog");
const storyStage = document.querySelector("#story-stage");
const storyProgress = document.querySelector("#story-progress");
const previousButton = document.querySelector("#story-previous");
const nextButton = document.querySelector("#story-next");
const closeButton = document.querySelector("#story-close");
let activeFilter = "all";
let storyEvents = events;
let storyIndex = 0;

function element(tag, className, text) {
  const node = document.createElement(tag);
  if (className) node.className = className;
  if (text !== undefined) node.textContent = text;
  return node;
}

function statusClass(event) {
  if (event.kind === "failure") return event.status === "discarded" ? "discarded" : "failed";
  if (event.kind === "result") return "result";
  if (event.kind === "infrastructure") return "infrastructure";
  return "pivot";
}

function fact(label, value, stat = false) {
  const row = element("div", "event-fact");
  row.append(element("b", "", label));
  row.append(element("span", stat ? "event-stat" : "", value));
  return row;
}

function eventCard(event) {
  const card = element("article", "timeline-event");
  card.dataset.kind = event.kind;
  const time = element("time", "", `${event.date} · ${event.time}`);
  const meta = element("div", "meta-row");
  meta.append(element("span", `chip ${statusClass(event)}`, event.status));
  const title = element("h2", "", event.title);
  const summary = element("p", "", event.summary);
  const facts = element("div", "event-facts");
  if (event.finding) facts.append(fact("Finding", event.finding));
  if (event.pivot) facts.append(fact("Next", event.pivot));
  const measure = event.stat ? element("p", "event-stat", event.stat) : null;
  const details = element("details", "event-details");
  details.append(element("summary", "", "Finding and next step"), facts);
  const link = element("a", "event-link", "Source record →");
  link.href = event.href;
  card.append(time, meta, title, summary);
  if (measure) card.append(measure);
  card.append(details, link);
  return card;
}

function visibleEvents() {
  return activeFilter === "all" ? events : events.filter((event) => event.kind === activeFilter);
}

function renderTimeline() {
  timeline.replaceChildren();
  const selected = visibleEvents();
  if (!selected.length) {
    timeline.append(element("p", "timeline-empty", "No events match this filter."));
    return;
  }
  selected.forEach((event) => timeline.append(eventCard(event)));
}

function renderStory() {
  const event = storyEvents[storyIndex];
  if (!event) return;
  storyStage.replaceChildren();
  const time = element("time", "", `${event.date} · ${event.time}`);
  const meta = element("div", "meta-row");
  meta.append(element("span", `chip ${statusClass(event)}`, event.status));
  const title = element("h2", "", event.title);
  const summary = element("p", "story-summary", event.summary);
  const facts = element("div", "event-facts");
  if (event.finding) facts.append(fact("Finding", event.finding));
  if (event.pivot) facts.append(fact("Next", event.pivot));
  if (event.stat) facts.append(fact("Measure", event.stat, true));
  storyStage.append(time, meta, title, summary, facts);
  storyProgress.textContent = `${storyIndex + 1} / ${storyEvents.length}`;
  previousButton.disabled = storyIndex === 0;
  nextButton.disabled = storyIndex === storyEvents.length - 1;
  nextButton.textContent = storyIndex === storyEvents.length - 1 ? "End" : "Next →";
}

function openStory() {
  storyEvents = visibleEvents();
  storyIndex = 0;
  renderStory();
  if (typeof dialog.showModal === "function") dialog.showModal();
  else dialog.setAttribute("open", "");
}

function closeStory() {
  if (typeof dialog.close === "function") dialog.close();
  else dialog.removeAttribute("open");
}

function moveStory(delta) {
  const next = storyIndex + delta;
  if (next < 0) return;
  if (next >= storyEvents.length) {
    closeStory();
    return;
  }
  storyIndex = next;
  renderStory();
}

filters.forEach((button) => {
  button.addEventListener("click", () => {
    activeFilter = button.dataset.filter;
    filters.forEach((candidate) => candidate.classList.toggle("active", candidate === button));
    renderTimeline();
  });
});

presentButton.addEventListener("click", openStory);
previousButton.addEventListener("click", () => moveStory(-1));
nextButton.addEventListener("click", () => moveStory(1));
closeButton.addEventListener("click", closeStory);
dialog.addEventListener("click", (event) => {
  if (event.target === dialog) closeStory();
});

document.addEventListener("keydown", (event) => {
  if (!dialog.hasAttribute("open")) return;
  if (["ArrowRight", "PageDown", " "].includes(event.key)) {
    event.preventDefault();
    moveStory(1);
  } else if (["ArrowLeft", "PageUp"].includes(event.key)) {
    event.preventDefault();
    moveStory(-1);
  } else if (event.key === "Home") {
    storyIndex = 0;
    renderStory();
  } else if (event.key === "End") {
    storyIndex = storyEvents.length - 1;
    renderStory();
  }
});

renderTimeline();
if (new URLSearchParams(window.location.search).get("present") === "1") openStory();
