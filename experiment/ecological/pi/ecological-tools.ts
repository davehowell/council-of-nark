/**
 * Narrow ecological respondent tools.
 *
 * The extension never receives a source or controller path. It forwards strict
 * JSON frames over inherited pipes to the trusted Go mediator. Launchers must
 * map requests to fd 3 and responses to fd 4, disable every built-in tool, and
 * retain Pi's JSON event stream alongside the mediator transcript.
 */

import { randomUUID } from "node:crypto";
import { createReadStream, writeSync } from "node:fs";
import type { ExtensionAPI } from "@earendil-works/pi-coding-agent";
import { StringEnum } from "@earendil-works/pi-ai";
import { Type } from "typebox";

const ACTIVE_TOOLS = [
	"source_list",
	"source_read",
	"source_search",
	"run_focused_test",
	"submit_ecological_review",
];

type MediatorResponse = {
	request_id: string;
	sequence: number;
	ok: boolean;
	output?: string;
	error?: string;
	truncated: boolean;
	metadata?: Record<string, unknown>;
};

const pending = new Map<
	string,
	{ resolve: (response: MediatorResponse) => void; reject: (error: Error) => void }
>();
let inputBuffer = Buffer.alloc(0);
let channelError: Error | undefined;
let auditSequence = 0;
let providerTurns = 0;
let observedTokens = 0;
let submissionCalls = 0;

const auditEnabled = process.env.COUNCIL_ECOLOGICAL_AUDIT === "1";
const maxProviderTurns = Number.parseInt(process.env.COUNCIL_ECOLOGICAL_MAX_PROVIDER_TURNS || "0", 10);
const maxTotalTokens = Number.parseInt(process.env.COUNCIL_ECOLOGICAL_MAX_TOTAL_TOKENS || "0", 10);

function audit(kind: string, data: Record<string, unknown> = {}): void {
	if (!auditEnabled) return;
	auditSequence += 1;
	writeSync(5, `${JSON.stringify({ schema_version: 1, sequence: auditSequence, kind, ...data })}\n`, undefined, "utf8");
}

function failChannel(error: Error): void {
	channelError = error;
	for (const waiter of pending.values()) waiter.reject(error);
	pending.clear();
}

function handleFrame(frame: Buffer): void {
	if (frame.length === 0) return;
	let response: MediatorResponse;
	try {
		response = JSON.parse(frame.toString("utf8")) as MediatorResponse;
	} catch (error) {
		failChannel(new Error(`invalid mediator response: ${String(error)}`));
		return;
	}
	const waiter = pending.get(response.request_id);
	if (!waiter) {
		failChannel(new Error(`unexpected mediator response id: ${response.request_id || "<empty>"}`));
		return;
	}
	pending.delete(response.request_id);
	waiter.resolve(response);
}

const responses = createReadStream("", { fd: 4, autoClose: false });
responses.on("data", (chunk: Buffer | string) => {
	const bytes = typeof chunk === "string" ? Buffer.from(chunk) : chunk;
	inputBuffer = Buffer.concat([inputBuffer, bytes]);
	for (;;) {
		const newline = inputBuffer.indexOf(0x0a);
		if (newline < 0) break;
		const frame = inputBuffer.subarray(0, newline);
		inputBuffer = inputBuffer.subarray(newline + 1);
		handleFrame(frame);
	}
});
responses.on("error", (error) => failChannel(error));
responses.on("end", () => failChannel(new Error("mediator response channel closed")));

async function mediate(requestId: string, tool: string, args: unknown): Promise<MediatorResponse> {
	if (channelError) throw channelError;
	if (!requestId) throw new Error("Pi did not provide a tool-call id");
	if (pending.has(requestId)) throw new Error(`duplicate tool-call id: ${requestId}`);
	const promise = new Promise<MediatorResponse>((resolve, reject) => {
		pending.set(requestId, { resolve, reject });
	});
	try {
		const frame = `${JSON.stringify({ request_id: requestId, tool, arguments: args })}\n`;
		writeSync(3, frame, undefined, "utf8");
	} catch (error) {
		pending.delete(requestId);
		throw error;
	}
	return promise;
}

function registerMediatedTool(
	pi: ExtensionAPI,
	definition: {
		name: string;
		label: string;
		description: string;
		parameters: ReturnType<typeof Type.Object>;
	},
): void {
	pi.registerTool({
		...definition,
		async execute(toolCallId, params) {
			const response = await mediate(toolCallId, definition.name, params);
			if (!response.ok) throw new Error(response.error || "mediator rejected the request");
			return {
				content: [{ type: "text", text: response.output || "" }],
				details: {
					sequence: response.sequence,
					truncated: response.truncated,
					...(response.metadata || {}),
				},
			};
		},
	});
}

export default function ecologicalTools(pi: ExtensionAPI): void {
	registerMediatedTool(pi, {
		name: "source_list",
		label: "List source",
		description:
			"List files and directories beneath a relative path in the sanitized source snapshot. No Git or sibling paths are available. Depth is capped at 3 and output is budgeted.",
		parameters: Type.Object({
			path: Type.Optional(Type.String({ description: "Relative directory path; defaults to ." })),
			depth: Type.Optional(Type.Integer({ minimum: 1, maximum: 3, description: "Recursive depth; defaults to 1" })),
		}),
	});

	registerMediatedTool(pi, {
		name: "source_read",
		label: "Read source",
		description:
			"Read numbered lines from one UTF-8 text file in the sanitized source snapshot. Paths are relative. Binary files, symlinks, traversal, and files over 2 MiB are denied. At most 400 lines are returned.",
		parameters: Type.Object({
			path: Type.String({ description: "Relative source-file path" }),
			start_line: Type.Optional(Type.Integer({ minimum: 1, description: "First line, 1-indexed; defaults to 1" })),
			line_count: Type.Optional(Type.Integer({ minimum: 1, maximum: 400, description: "Lines to return; defaults to 400" })),
		}),
	});

	registerMediatedTool(pi, {
		name: "source_search",
		label: "Search source",
		description:
			"Search UTF-8 source lines with a Go RE2 regular expression. Restrict by relative directory and optional basename glob when useful. Results are capped at 100 matches; there is no shell, Git, or internet search.",
		parameters: Type.Object({
			pattern: Type.String({ description: "Go RE2 regular expression, at most 512 UTF-8 bytes" }),
			path: Type.Optional(Type.String({ description: "Relative directory path; defaults to ." })),
			include: Type.Optional(Type.String({ description: "Optional basename glob such as *.go" })),
		}),
	});

	registerMediatedTool(pi, {
		name: "run_focused_test",
		label: "Run focused test",
		description:
			"Run one predeclared regression target in a separate network-denied sandbox with the frozen toolchain and dependency closure. The hidden test source, controller paths, and shell are not exposed. At most 3 calls are allowed.",
		parameters: Type.Object({
			target: StringEnum(["unicode-tokenizer-regression"] as const, {
				description: "The only preregistered test target",
			}),
		}),
	});

	pi.registerTool({
		name: "submit_ecological_review",
		label: "Submit ecological review",
		description:
			"Submit the final structured review. Use this exactly once as the final action. Report only supported findings; do not invent findings to fill the allowance.",
		parameters: Type.Object({
			review_summary: Type.String({ description: "Concise overall diagnosis" }),
			findings: Type.Array(
				Type.Object({
					title: Type.String(),
					locations: Type.Array(
						Type.Object({
							path: Type.String({ description: "Relative source path" }),
							line_start: Type.Integer({ minimum: 1 }),
							line_end: Type.Integer({ minimum: 1 }),
						}),
						{ minItems: 1, maxItems: 4 },
					),
					claim: Type.String({ description: "Unsafe assumption or defect" }),
					mechanism: Type.String({ description: "Why and under what conditions it occurs" }),
					consequence: Type.String(),
					correction: Type.String({ description: "Correct, scoped remedy" }),
					allocation_and_scope: Type.String({ description: "Hot-path allocation and compatibility implications" }),
					regression_tests: Type.Array(
						Type.Object({
							case: Type.String(),
							expected: Type.String(),
							why_discriminating: Type.String(),
						}),
						{ maxItems: 12 },
					),
					evidence: Type.Array(Type.String(), { minItems: 1, maxItems: 12 }),
					confidence: StringEnum(["high", "medium", "low"] as const),
				}),
				{ minItems: 1, maxItems: 8 },
			),
			uncertainties: Type.Array(Type.String(), { maxItems: 8 }),
		}),
		async execute(toolCallId, params) {
			audit("final_submission", { tool_call_id: toolCallId, details: params });
			return {
				content: [{ type: "text", text: `Submitted ${params.findings.length} finding(s).` }],
				details: params,
				terminate: true,
			};
		},
	});

	if (process.env.COUNCIL_ECOLOGICAL_DOCTOR === "1") {
		pi.registerCommand("ecological-mediator-health", {
			description: "Controller-only inherited-pipe health check; makes no model call",
			handler: async (_args, ctx) => {
				const response = await mediate(`health-${randomUUID()}`, "source_list", { path: ".", depth: 1 });
				if (!response.ok) throw new Error(response.error || "mediator health check failed");
				ctx.ui.notify(
					JSON.stringify({ ecological_mediator_health: true, sequence: response.sequence }),
					"info",
				);
			},
		});
	}

	pi.on("before_provider_request", (event, ctx) => {
		if (
			(maxProviderTurns > 0 && providerTurns >= maxProviderTurns) ||
			(maxTotalTokens > 0 && observedTokens >= maxTotalTokens)
		) {
			audit("budget_stop", { provider_turns: providerTurns, observed_tokens: observedTokens });
			ctx.abort();
			throw new Error("ecological aggregate provider budget exhausted");
		}
		providerTurns += 1;
		audit("provider_request", { provider_turn: providerTurns, payload: event.payload });
	});

	pi.on("after_provider_response", (event) => {
		audit("provider_response", {
			provider_turn: providerTurns,
			status: event.status,
			headers: event.headers,
		});
	});

	pi.on("message_end", (event) => {
		if (event.message.role !== "assistant") return;
		const usage = event.message.usage;
		const total = usage.totalTokens || usage.input + usage.output + usage.cacheRead + usage.cacheWrite;
		observedTokens += total;
		audit("assistant_usage", {
			provider_turn: providerTurns,
			usage,
			observed_tokens: observedTokens,
		});
	});

	pi.on("tool_call", (event) => {
		if (submissionCalls > 0) {
			audit("tool_after_submission", { tool_call_id: event.toolCallId, tool_name: event.toolName });
			return { block: true, reason: "the final ecological submission has already been made", terminate: true };
		}
		if (event.toolName === "submit_ecological_review") {
			submissionCalls += 1;
		}
	});

	pi.on("session_start", () => {
		pi.setActiveTools(ACTIVE_TOOLS);
		audit("session_start", { active_tools: pi.getActiveTools() });
	});
}
