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

    useEffect(() => {
        async function loadGraph() {
            try {
                const response =
                    await api.get("/graph");

                const graph =
                    response.data;

                const flowNodes: Node[] =
                    graph.nodes.map(
                        (
                            node: any,
                            index: number,
                        ) => ({
                            id: node.id,
                            data: {
                                label: (
                                    <div
                                        style={{
                                            padding: "10px 16px",
                                            borderRadius: "12px",
                                            background: "#18181b",
                                            color: "white",
                                            border: "1px solid #27272a",
                                            fontSize: "14px",
                                            fontWeight: 500,
                                            minWidth: "120px",
                                            textAlign: "center",
                                        }}
                                    >
                                        {node.id}
                                    </div>
                                ),
                            },
                            position: {
                                x:
                                    (index % 5) * 250,
                                y:
                                    Math.floor(
                                        index / 5
                                    ) * 150,
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
                            label:
                                edge.label,
                            animated: true,
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
        flex
        h-screen
        bg-zinc-950
        "
        >
            <div
                style={{
                    position: "absolute",
                    top: 20,
                    left: 20,
                    zIndex: 1000,
                    background: "#18181b",
                    color: "white",
                    padding: "12px 18px",
                    borderRadius: "12px",
                    border: "1px solid #27272a",
                }}
            >
                <div
                    style={{
                        fontSize: "18px",
                        fontWeight: "bold",
                    }}
                >
                    DevBrain
                </div>

                <div
                    style={{
                        fontSize: "12px",
                        color: "#a1a1aa",
                    }}
                >
                    Interactive Knowledge Graph
                </div>
            </div>
            <div
                className="
        flex-1
        "
            >
                <ReactFlow
                    nodes={nodes}
                    edges={edges}
                    fitView
                    fitViewOptions={{
                        padding: 0.3,
                    }}
                    onNodeClick={(
                        _,
                        node,
                    ) => {
                        loadConcept(
                            node.id
                        );
                    }}
                >
                    <Background />
                    <Controls />
                </ReactFlow>
            </div>

            {selectedConcept && (
                <div
                    className="
        w-96
        border-l
        border-zinc-800
        bg-zinc-900
        text-white
        p-5
        overflow-auto
        "
                >

                    <h2
                        className="
            text-2xl
            font-bold
            mb-4
            "
                    >
                        {
                            selectedConcept.concept
                        }
                    </h2>

                    <div
                        className="
            flex
            flex-col
            gap-3
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
    border-zinc-700
    bg-zinc-800
    rounded-xl
    p-4
    "
                                >

                                    <div>
                                        <strong>
                                            {
                                                rel.source
                                            }
                                        </strong>
                                    </div>

                                    <div>
                                        <span
                                            className="
        inline-block
        px-2
        py-1
        rounded-full
        bg-blue-500/20
        text-blue-300
        text-xs
        font-medium
        "
                                        >
                                            {rel.relation}
                                        </span>
                                    </div>

                                    <div>
                                        <strong>
                                            {
                                                rel.target
                                            }
                                        </strong>
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