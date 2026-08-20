from __future__ import annotations

import argparse
from itertools import permutations
import random
from typing import Any

from .common import load_json, opaque_id, resolve_run, utc_now, write_json

ALL_ROLES = [
    "simplicity",
    "correctness",
    "architecture",
    "security",
    "cost",
    "entropy",
    "language",
]


class Builder:
    def __init__(self, config: dict[str, Any]):
        self.config = config
        self.seed = config["seed"]
        self.calls: list[dict[str, Any]] = []
        self.sets: list[dict[str, Any]] = []

    def call(
        self,
        semantic: tuple[object, ...],
        *,
        packet: str,
        provider: dict[str, Any],
        prompt_spec: dict[str, Any],
        depends_on: list[str] | None = None,
        phase: str = "review",
    ) -> dict[str, Any]:
        call_id = "c-" + opaque_id(self.seed, "call", *semantic)
        value = {
            "call_id": call_id,
            "reviewer_id": "r-" + opaque_id(self.seed, "reviewer", *semantic, length=10),
            "packet": packet,
            "provider": provider,
            "prompt_spec": prompt_spec,
            "depends_on": depends_on or [],
            "phase": phase,
            "semantic": list(semantic),
        }
        self.calls.append(value)
        return value

    def output_set(
        self,
        semantic: tuple[object, ...],
        *,
        packet: str,
        call_ids: list[str],
        cost_call_ids: list[str],
        metadata: dict[str, Any],
    ) -> None:
        self.sets.append(
            {
                "set_id": "s-" + opaque_id(self.seed, "set", *semantic),
                "packet": packet,
                "call_ids": call_ids,
                "cost_call_ids": cost_call_ids,
                "metadata": metadata,
                "semantic": list(semantic),
            }
        )


def stage_a(builder: Builder) -> None:
    config = builder.config
    provider = config["provider"]
    for packet in config["packets"]:
        for repeat in range(config["repetitions"]):
            for arm in config["arms"]:
                prefix = ("stage_a", packet, repeat, arm)
                if arm in {"S0", "S1", "S2"}:
                    spec = {
                        "S0": {"kind": "generic"},
                        "S1": {"kind": "omnibus", "wrapper": "functional"},
                        "S2": {"kind": "omnibus", "wrapper": "fictional"},
                    }[arm]
                    review = builder.call(prefix + ("review",), packet=packet, provider=provider, prompt_spec=spec)
                    builder.output_set(
                        prefix + ("final",),
                        packet=packet,
                        call_ids=[review["call_id"]],
                        cost_call_ids=[review["call_id"]],
                        metadata={"design": "stage_a", "arm": arm, "repeat": repeat, "kind": "final"},
                    )
                    continue

                reviews: list[dict[str, Any]] = []
                for index, role in enumerate(ALL_ROLES):
                    if arm == "M0":
                        spec = {"kind": "omnibus", "wrapper": "functional", "sample": index}
                    elif arm == "M1":
                        spec = {"kind": "specialist", "wrapper": "functional", "role": role}
                    elif arm == "M2":
                        spec = {"kind": "specialist", "wrapper": "fictional", "role": role}
                    else:
                        raise ValueError(f"Unknown Stage A arm: {arm}")
                    reviews.append(
                        builder.call(prefix + ("review", index, role), packet=packet, provider=provider, prompt_spec=spec)
                    )
                review_ids = [review["call_id"] for review in reviews]
                fused = builder.call(
                    prefix + ("fuse",),
                    packet=packet,
                    provider=provider,
                    prompt_spec={"kind": "fuser"},
                    depends_on=review_ids,
                    phase="fuse",
                )
                builder.output_set(
                    prefix + ("raw_union",),
                    packet=packet,
                    call_ids=review_ids,
                    cost_call_ids=review_ids,
                    metadata={"design": "stage_a", "arm": arm, "repeat": repeat, "kind": "raw_union"},
                )
                builder.output_set(
                    prefix + ("fused",),
                    packet=packet,
                    call_ids=[fused["call_id"]],
                    cost_call_ids=review_ids + [fused["call_id"]],
                    metadata={"design": "stage_a", "arm": arm, "repeat": repeat, "kind": "fused"},
                )


def topology(builder: Builder) -> None:
    config = builder.config
    provider = config["provider"]
    roles = config["roles"]
    orders = list(permutations(roles)) if config.get("all_role_orders") else [tuple(roles)]
    wrapper = config["wrapper"]
    for packet in config["packets"]:
        for repeat in range(config["repetitions"]):
            for order in orders:
                order_name = "-".join(order)
                for topology_name in config["topologies"]:
                    prefix = ("topology", packet, repeat, order_name, topology_name)
                    reviews: list[dict[str, Any]] = []
                    previous: list[str] = []
                    for hop, role in enumerate(order):
                        dynamic = topology_name == "chain" and hop > 0
                        spec = {
                            "kind": "chain" if dynamic else "specialist",
                            "wrapper": wrapper,
                            "role": role,
                            "hop": hop + 1,
                        }
                        dependencies = previous[-1:] if dynamic else []
                        review = builder.call(
                            prefix + ("review", hop, role),
                            packet=packet,
                            provider=provider,
                            prompt_spec=spec,
                            depends_on=dependencies,
                        )
                        reviews.append(review)
                        previous.append(review["call_id"])
                        if topology_name == "chain":
                            builder.output_set(
                                prefix + (f"hop-{hop + 1}",),
                                packet=packet,
                                call_ids=[review["call_id"]],
                                cost_call_ids=previous.copy(),
                                metadata={
                                    "design": "topology",
                                    "topology": topology_name,
                                    "order": order_name,
                                    "repeat": repeat,
                                    "kind": f"hop-{hop + 1}",
                                },
                            )
                    review_ids = [review["call_id"] for review in reviews]
                    fused = builder.call(
                        prefix + ("fuse",),
                        packet=packet,
                        provider=provider,
                        prompt_spec={"kind": "fuser"},
                        depends_on=review_ids,
                        phase="fuse",
                    )
                    if topology_name == "fanout":
                        builder.output_set(
                            prefix + ("raw_union",),
                            packet=packet,
                            call_ids=review_ids,
                            cost_call_ids=review_ids,
                            metadata={
                                "design": "topology",
                                "topology": topology_name,
                                "order": order_name,
                                "repeat": repeat,
                                "kind": "raw_union",
                            },
                        )
                    builder.output_set(
                        prefix + ("fused",),
                        packet=packet,
                        call_ids=[fused["call_id"]],
                        cost_call_ids=review_ids + [fused["call_id"]],
                        metadata={
                            "design": "topology",
                            "topology": topology_name,
                            "order": order_name,
                            "repeat": repeat,
                            "kind": "fused",
                        },
                    )


def provider_pair(builder: Builder) -> None:
    config = builder.config
    for packet in config["packets"]:
        for repeat in range(config["repetitions"]):
            for provider_index, provider in enumerate(config["providers"]):
                for wrapper in config["wrappers"]:
                    prefix = ("provider_pair", packet, repeat, provider_index, wrapper)
                    review = builder.call(
                        prefix + ("review",),
                        packet=packet,
                        provider=provider,
                        prompt_spec={"kind": "specialist", "wrapper": wrapper, "role": config["role"]},
                    )
                    builder.output_set(
                        prefix + ("final",),
                        packet=packet,
                        call_ids=[review["call_id"]],
                        cost_call_ids=[review["call_id"]],
                        metadata={
                            "design": "provider_pair",
                            "provider_index": provider_index,
                            "adapter": provider["adapter"],
                            "model": provider["model"],
                            "wrapper": wrapper,
                            "repeat": repeat,
                            "kind": "final",
                        },
                    )


def build(config: dict[str, Any]) -> dict[str, Any]:
    builder = Builder(config)
    design = config["design"]
    if design == "stage_a":
        stage_a(builder)
    elif design == "topology":
        topology(builder)
    elif design == "provider_pair":
        provider_pair(builder)
    else:
        raise ValueError(f"Unknown design: {design}")
    random.Random(config["seed"]).shuffle(builder.calls)
    return {
        "schema_version": 1,
        "created_at": utc_now(),
        "design": design,
        "seed": config["seed"],
        "calls": builder.calls,
        "output_sets": builder.sets,
        "counts": {"calls": len(builder.calls), "output_sets": len(builder.sets)},
    }


def main() -> int:
    parser = argparse.ArgumentParser(description="Create a deterministic call plan for a frozen run")
    parser.add_argument("run")
    args = parser.parse_args()
    run = resolve_run(args.run)
    if (run / "plan.json").exists():
        raise SystemExit("plan.json already exists; do not replace a frozen run plan")
    config = load_json(run / "config.json")
    plan = build(config)
    write_json(run / "plan.json", plan)
    print(f"Planned {plan['counts']['calls']} calls and {plan['counts']['output_sets']} output sets.")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
