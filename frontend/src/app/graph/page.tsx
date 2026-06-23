"use client";

import { useEffect, useState } from "react";
import { api } from "@/lib/api";
import { getLayoutedElements } from "@/lib/layout";

import ReactFlow, {
    Background,
    Controls,
    Node,
    Edge,
} from "reactflow";

import "reactflow/dist/style.css";

export default function GraphPage() {
    const [nodes, setNodes] = useState<Node[]>([]);
    const [edges, setEdges] = useState<Edge[]>([]);
    const [selectedConcept, setSelectedConcept] =
        useState<any>(null);

    const [search, setSearch] =
        useState("");

    const [reactFlowInstance, setReactFlowInstance] =
        useState<any>(null);

    async function loadConcept(
        name: string,
    ) {
        try {
            const response =
                await api.get(
                    `/concepts/${encodeURIComponent(name)}`
                );

            setSelectedConcept(
                response.data
            );

        } catch (err) {
            console.error(err);
        }
    }

    function searchNode() {
        const found = nodes.find(
            (n) =>
                n.id.toLowerCase() ===
                search.toLowerCase()
        );

        if (!found) {
            alert("Concept not found");
            return;
        }

        loadConcept(found.id);

        reactFlowInstance?.setCenter(
            found.position.x,
            found.position.y,
            {
                zoom: 1.5,
                duration: 800,
            }
        );
    }

    useEffect(() => {
        async function loadGraph() {
            try {
                const response =
                    await api.get("/graph");

                const graph =
                    response.data;

                const flowNodes: Node[] =
                    graph.nodes.map(
                        (node: any) => ({
                            id: node.id,

                            data: {
                                label: node.id,
                            },

                            position: {
                                x: 0,
                                y: 0,
                            },

                            style: {
                                background: "#18181b",
                                color: "#ffffff",
                                border:
                                    "1px solid #27272a",
                                borderRadius: "16px",
                                minWidth: "180px",
                                padding: "10px",
                                textAlign: "center",
                                fontWeight: 600,
                                fontSize: "14px",
                            },
                        })
                    );

                const flowEdges: Edge[] =
                    graph.edges.map(
                        (
                            edge: any,
                            index: number,
                        ) => ({
                            id:
                                `${edge.source}-${edge.target}-${index}`,

                            source:
                                edge.source,

                            target:
                                edge.target,

                            label: "",

                            animated: false,

                            style: {
                                stroke:
                                    "#52525b",
                            },
                        })
                    );

                const layouted =
                    getLayoutedElements(
                        flowNodes,
                        flowEdges
                    );

                setNodes(
                    layouted.nodes
                );

                setEdges(
                    layouted.edges
                );

            } catch (err) {
                console.error(err);
            }
        }

        loadGraph();
    }, []);

    return (
        <div
            className="
            h-screen
            w-screen
            bg-zinc-950
            relative
            "
        >

            {/* DevBrain Card */}

            <div
                className="
                absolute
                top-6
                left-6
                z-50

                bg-zinc-900/90
                backdrop-blur-xl

                border
                border-zinc-800

                rounded-3xl

                p-6

                shadow-2xl
                "
            >
                <h1
                    className="
                    text-3xl
                    font-bold
                    text-white
                    "
                >
                    DevBrain
                </h1>

                <p
                    className="
                    text-zinc-400
                    "
                >
                    Interactive Knowledge Graph
                </p>

                <div
                    className="
                    mt-4
                    flex
                    gap-6
                    "
                >
                    <div>
                        <div
                            className="
                            text-white
                            font-bold
                            text-xl
                            "
                        >
                            {nodes.length}
                        </div>

                        <div
                            className="
                            text-zinc-500
                            text-sm
                            "
                        >
                            Concepts
                        </div>
                    </div>

                    <div>
                        <div
                            className="
                            text-white
                            font-bold
                            text-xl
                            "
                        >
                            {edges.length}
                        </div>

                        <div
                            className="
                            text-zinc-500
                            text-sm
                            "
                        >
                            Relations
                        </div>
                    </div>
                </div>
            </div>

            {/* Search */}

            <div
                className="
                absolute
                top-6
                left-1/2
                -translate-x-1/2

                z-50

                flex
                gap-3
                "
            >
                <input
                    value={search}
                    onChange={(e) =>
                        setSearch(
                            e.target.value
                        )
                    }
                    placeholder="Search concept..."
                    className="
                    w-96

                    bg-zinc-900/90
                    backdrop-blur-xl

                    border
                    border-zinc-800

                    rounded-2xl

                    px-5
                    py-3

                    text-white

                    outline-none
                    "
                />

                <button
                    onClick={searchNode}
                    className="
                    px-5
                    py-3

                    rounded-2xl

                    bg-blue-600
                    text-white

                    font-medium
                    "
                >
                    Search
                </button>
            </div>

            <ReactFlow
                nodes={nodes}
                edges={edges}
                onInit={setReactFlowInstance}
                fitView
                fitViewOptions={{
                    padding: 0.3,
                }}
                onPaneClick={() =>
                    setSelectedConcept(null)
                }
                onNodeClick={(
                    _,
                    node,
                ) => {
                    loadConcept(
                        node.id
                    );
                }}
            >
                <Background
                    gap={24}
                    size={1}
                    color="#27272a"
                />

                <Controls />
            </ReactFlow>

            {selectedConcept && (
                <div
                    className="
                    absolute

                    top-6
                    right-6
                    bottom-6

                    w-[420px]

                    z-50

                    bg-zinc-900/90
                    backdrop-blur-xl

                    border
                    border-zinc-800

                    rounded-3xl

                    p-6

                    overflow-y-auto

                    shadow-2xl
                    "
                >
                    <div
                        className="
  flex
  items-center
  justify-between
  mb-6
  "
                    >
                        <h2
                            className="
    text-3xl
    font-bold
    text-white
    "
                        >
                            {selectedConcept.concept}
                        </h2>

                        <button
                            onClick={() =>
                                setSelectedConcept(null)
                            }
                            className="
  h-10
  w-10

  rounded-full

  bg-zinc-800

  text-zinc-400

  hover:bg-zinc-700
  hover:text-white

  transition
  "
                        >
                            ✕
                        </button>
                    </div>

                    <div
                        className="
                        flex
                        flex-col
                        gap-4
                        "
                    >
                        {selectedConcept.relationships.map(
                            (
                                rel: any,
                                index: number,
                            ) => (
                                <div
                                    key={index}
                                    className="
                                    border
                                    border-zinc-800

                                    bg-zinc-900

                                    rounded-2xl

                                    p-4
                                    "
                                >
                                    <div
                                        className="
                                        text-white
                                        font-semibold
                                        "
                                    >
                                        {
                                            rel.source
                                        }
                                    </div>

                                    <div
                                        className="
                                        mt-2
                                        mb-2
                                        "
                                    >
                                        <span
                                            className="
                                            inline-block

                                            px-3
                                            py-1

                                            rounded-full

                                            bg-blue-500/20
                                            text-blue-300

                                            text-xs
                                            "
                                        >
                                            {
                                                rel.relation
                                            }
                                        </span>
                                    </div>

                                    <div
                                        className="
                                        text-white
                                        "
                                    >
                                        {
                                            rel.target
                                        }
                                    </div>
                                </div>
                            )
                        )}
                    </div>
                </div>
            )}
        </div>
    );
}