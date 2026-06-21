"use client";

import { useEffect, useState } from "react";
import { api } from "@/lib/api";

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

  useEffect(() => {
    async function loadGraph() {
      try {
        const response = await api.get("/graph");

        const graph = response.data;

        const flowNodes: Node[] = graph.nodes.map(
          (node: any, index: number) => ({
            id: node.id,
            data: {
              label: node.id,
            },
            position: {
              x: (index % 5) * 250,
              y: Math.floor(index / 5) * 150,
            },
          })
        );

        const flowEdges: Edge[] = graph.edges.map(
          (edge: any, index: number) => ({
            id: `${edge.source}-${edge.target}-${index}`,
            source: edge.source,
            target: edge.target,
            label: edge.label,
            animated: true,
          })
        );

        setNodes(flowNodes);
        setEdges(flowEdges);
      } catch (err) {
        console.error(err);
      }
    }

    loadGraph();
  }, []);

  return (
    <div style={{ width: "100vw", height: "100vh" }}>
      <ReactFlow
        nodes={nodes}
        edges={edges}
        fitView
      >
        <Background />
        <Controls />
      </ReactFlow>
    </div>
  );
}