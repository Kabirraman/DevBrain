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
    const [selectedNode, setSelectedNode] =
        useState<string | null>(null);
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

        setSelectedNode(found.id);
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
                                background: "rgba(24,24,27,.85)",
                                color: "#ffffff",
                                border: "1px solid #27272a",
                                borderRadius: "18px",
                                minWidth: "160px",
                                padding: "8px",

                                boxShadow:
                                    "0 8px 30px rgba(0,0,0,.35)",
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
                                stroke: "#3f3f46",
                                strokeWidth: 1,
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

    const highlightedNodes =
        nodes.map((node) => {

            if (!selectedNode)
                return node;

            const connected =
                edges.some(
                    (e) =>
                        (e.source === selectedNode &&
                            e.target === node.id) ||

                        (e.target === selectedNode &&
                            e.source === node.id)
                );

            const selected =
                node.id === selectedNode;

            return {
                ...node,

                style: {
                    ...node.style,

                    opacity:
                        selected || connected
                            ? 1
                            : 0.15,

                    border:
                        selected
                            ? "1px solid #3b82f6"
                            : "1px solid #27272a",

                    boxShadow:
                        selected
                            ? "0 0 25px rgba(59,130,246,.45)"
                            : "0 8px 30px rgba(0,0,0,.35)",
                },
            };
        });

    const highlightedEdges =
        edges.map((edge) => {

            const active =
                edge.source === selectedNode ||
                edge.target === selectedNode;

            return {
                ...edge,

                style: {
                    stroke: "#3f3f46",

                    opacity:
                        !selectedNode
                            ? 0.35
                            : active
                                ? 1
                                : 0.05,

                    strokeWidth:
                        active
                            ? 2
                            : 1,
                },
            };
        });

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

                rounded-2xl

                p-6

                shadow-2xl
                "
            >
                <h1
    className="
    text-2xl
    font-semibold
    tracking-tight
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
                    Knowledge Graph Explorer
                </p>

                <div
                    className="
    mt-5
    flex
    gap-3
    "
                >
                    <div
                        className="
        px-4
        py-3

        rounded-2xl

        bg-white/5

        border
        border-white/10
        "
                    >
                        <div
                            className="
            text-white
            font-bold
            text-lg
            "
                        >
                            {nodes.length}
                        </div>

                        <div
                            className="
            text-zinc-400
            text-xs
            "
                        >
                            Concepts
                        </div>
                    </div>

                    <div
                        className="
        px-4
        py-3

        rounded-2xl

        bg-white/5

        border
        border-white/10
        "
                    >
                        <div
                            className="
            text-white
            font-bold
            text-lg
            "
                        >
                            {edges.length}
                        </div>

                        <div
                            className="
            text-zinc-400
            text-xs
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
                    placeholder="Search concepts, tools, releases..."
                    className="
w-[420px]

h-14

bg-zinc-900/80
backdrop-blur-xl

border
border-white/10

rounded-full

px-6

text-white

shadow-xl

outline-none




                    
                    py-3


                    "
                />

                <button
                    onClick={searchNode}
                    className="
h-14

px-8

rounded-full

bg-blue-600

hover:bg-blue-500

transition

font-semibold

shadow-lg
"
                >
                    Search
                </button>
            </div>

            <div
                className={`
    h-full
    transition-all
    duration-300
    ${selectedConcept
                        ? "blur-[2px] opacity-60"
                        : ""
                    }
    `}
            >
                <ReactFlow
                    nodes={highlightedNodes}
                    edges={highlightedEdges}
                    onInit={setReactFlowInstance}
                    proOptions={{
                        hideAttribution: true,
                    }}
                    fitView
                    fitViewOptions={{
                        padding: 0.3,
                    }}
                    onPaneClick={() =>
                        setSelectedConcept(null)
                    }
                    onNodeClick={(_, node) => {

                        setSelectedNode(node.id);

                        loadConcept(node.id);
                    }}
                >
                    <Background
                       gap={24}
                       size={1}
                       color="#3f3f46"
                        />

                    <Controls />
                </ReactFlow>
            </div>

            {selectedConcept ? (
                <div
                    className="
        absolute
        top-6
        right-6
        bottom-6

        w-[520px]

        bg-zinc-900/80
        backdrop-blur-2xl

        border
        border-white/10

        rounded-2xl

        p-8

        shadow-[0_20px_80px_rgba(0,0,0,.6)]

        overflow-y-auto

        z-50
        "
                >
                    <div
                        className="
                        animate-in
slide-in-from-right
duration-300
  flex
  items-center
  justify-between
  mb-6
  "
                    >
                        <h2
                            className="
    text-2xl
    font-bold
    text-white
    "
                        >
                            {selectedConcept.concept}
                        </h2>

                        <button
                            onClick={() => {
                                setSelectedConcept(null);
                                setSelectedNode(null);
                            }}
                            className="
 h-10
w-10

flex
items-center
justify-center

rounded-full

bg-white/5

border
border-white/10

text-zinc-400

hover:bg-white/10
hover:text-white

transition
  "
                        >
                            ×
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
bg-white/[0.03]

border
border-white/10

rounded-2xl

p-4

hover:border-blue-500/30

transition
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

                                           bg-blue-500/10
text-blue-400
border
border-blue-500/20

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
            ) : (
//click any node...
    <div
        className="
        absolute

        bottom-6
        right-6

        bg-zinc-900/80
        backdrop-blur-xl

        border
        border-white/10

        rounded-2xl

        px-5
        py-4

        text-zinc-400
        text-sm

        z-50
        "
    >
        <div className="flex flex-col gap-2">

    <div className="text-white font-medium">
        Explore Concepts
    </div>

    <div className="text-zinc-400 text-sm">
        Select a node to inspect relationships
    </div>

</div>
    </div>

)}
        </div>
    );
}