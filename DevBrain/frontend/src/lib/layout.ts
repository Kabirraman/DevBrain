import dagre from "dagre";
import { Edge, Node } from "reactflow";

const dagreGraph = new dagre.graphlib.Graph();

dagreGraph.setDefaultEdgeLabel(() => ({}));

const nodeWidth = 220;
const nodeHeight = 60;

export function getLayoutedElements(
  nodes: Node[],
  edges: Edge[],
) {
  dagreGraph.setGraph({
  rankdir: "TB",
  nodesep: 40,
  ranksep: 80,
  marginx: 20,
  marginy: 20,
});

  nodes.forEach((node) => {
    dagreGraph.setNode(node.id, {
      width: nodeWidth,
      height: nodeHeight,
    });
  });

  edges.forEach((edge) => {
    dagreGraph.setEdge(
      edge.source,
      edge.target,
    );
  });

  dagre.layout(dagreGraph);

  nodes.forEach((node) => {
    const nodeWithPosition =
      dagreGraph.node(node.id);

    node.position = {
      x:
        nodeWithPosition.x -
        nodeWidth / 2,
      y:
        nodeWithPosition.y -
        nodeHeight / 2,
    };
  });

  return {
    nodes,
    edges,
  };
}