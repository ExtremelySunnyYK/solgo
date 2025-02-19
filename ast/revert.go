package ast

import (
	"github.com/goccy/go-json"

	v3 "github.com/cncf/xds/go/xds/type/v3"
	ast_pb "github.com/unpackdev/protos/dist/go/ast"
	"github.com/unpackdev/solgo/parser"
)

// RevertStatement represents a revert statement in the AST.
type RevertStatement struct {
	*ASTBuilder

	Id         int64            `json:"id"`         // Unique identifier for the RevertStatement node.
	NodeType   ast_pb.NodeType  `json:"node_type"`  // Type of the AST node.
	Src        SrcNode          `json:"src"`        // Source location information.
	Arguments  []Node[NodeType] `json:"arguments"`  // List of argument expressions.
	Expression Node[NodeType]   `json:"expression"` // Expression within the revert statement.
}

// NewRevertStatement creates a new RevertStatement node with a given ASTBuilder.
func NewRevertStatement(b *ASTBuilder) *RevertStatement {
	return &RevertStatement{
		ASTBuilder: b,

		Id:        b.GetNextID(),
		NodeType:  ast_pb.NodeType_REVERT_STATEMENT,
		Arguments: make([]Node[NodeType], 0),
	}
}

// SetReferenceDescriptor sets the reference descriptions of the RevertStatement node.
func (r *RevertStatement) SetReferenceDescriptor(refId int64, refDesc *TypeDescription) bool {
	return false
}

// GetId returns the ID of the RevertStatement node.
func (r *RevertStatement) GetId() int64 {
	return r.Id
}

// GetType returns the NodeType of the RevertStatement node.
func (r *RevertStatement) GetType() ast_pb.NodeType {
	return r.NodeType
}

// GetSrc returns the SrcNode of the RevertStatement node.
func (r *RevertStatement) GetSrc() SrcNode {
	return r.Src
}

// GetArguments returns the list of argument expressions.
func (r *RevertStatement) GetArguments() []Node[NodeType] {
	return r.Arguments
}

// GetExpression returns the expression within the revert statement.
func (r *RevertStatement) GetExpression() Node[NodeType] {
	return r.Expression
}

// GetNodes returns the child nodes of the RevertStatement node.
func (r *RevertStatement) GetNodes() []Node[NodeType] {
	toReturn := make([]Node[NodeType], 0)
	toReturn = append(toReturn, r.Arguments...)
	toReturn = append(toReturn, r.Expression)
	return toReturn
}

// MarshalJSON marshals the RevertStatement node into a JSON byte slice.
func (r *RevertStatement) UnmarshalJSON(data []byte) error {
	var tempMap map[string]json.RawMessage
	if err := json.Unmarshal(data, &tempMap); err != nil {
		return err
	}

	if id, ok := tempMap["id"]; ok {
		if err := json.Unmarshal(id, &r.Id); err != nil {
			return err
		}
	}

	if nodeType, ok := tempMap["node_type"]; ok {
		if err := json.Unmarshal(nodeType, &r.NodeType); err != nil {
			return err
		}
	}

	if src, ok := tempMap["src"]; ok {
		if err := json.Unmarshal(src, &r.Src); err != nil {
			return err
		}
	}

	if arguments, ok := tempMap["arguments"]; ok {
		var nodes []json.RawMessage
		if err := json.Unmarshal(arguments, &nodes); err != nil {
			return err
		}

		for _, tempNode := range nodes {
			var tempNodeMap map[string]json.RawMessage
			if err := json.Unmarshal(tempNode, &tempNodeMap); err != nil {
				return err
			}

			var tempNodeType ast_pb.NodeType
			if err := json.Unmarshal(tempNodeMap["node_type"], &tempNodeType); err != nil {
				return err
			}

			node, err := unmarshalNode(tempNode, tempNodeType)
			if err != nil {
				return err
			}
			r.Arguments = append(r.Arguments, node)
		}
	}

	if expression, ok := tempMap["expression"]; ok {
		if err := json.Unmarshal(expression, &r.Expression); err != nil {
			var tempNodeMap map[string]json.RawMessage
			if err := json.Unmarshal(expression, &tempNodeMap); err != nil {
				return err
			}

			var tempNodeType ast_pb.NodeType
			if err := json.Unmarshal(tempNodeMap["node_type"], &tempNodeType); err != nil {
				return err
			}

			node, err := unmarshalNode(expression, tempNodeType)
			if err != nil {
				return err
			}
			r.Expression = node
		}
	}

	return nil
}

// ToProto returns a protobuf representation of the RevertStatement node.
func (r *RevertStatement) ToProto() NodeType {
	proto := ast_pb.Revert{
		Id:         r.Id,
		NodeType:   r.NodeType,
		Src:        r.Src.ToProto(),
		Arguments:  make([]*v3.TypedStruct, 0),
		Expression: r.Expression.ToProto().(*v3.TypedStruct),
	}

	for _, arg := range r.Arguments {
		proto.Arguments = append(proto.Arguments, arg.ToProto().(*v3.TypedStruct))
	}

	return NewTypedStruct(&proto, "Revert")
}

// GetTypeDescription returns the TypeDescription of the RevertStatement node.
func (r *RevertStatement) GetTypeDescription() *TypeDescription {
	return &TypeDescription{
		TypeString:     "revert",
		TypeIdentifier: "$_t_revert",
	}
}

// Parse parses a revert statement context into the RevertStatement node.
func (r *RevertStatement) Parse(
	unit *SourceUnit[Node[ast_pb.SourceUnit]],
	contractNode Node[NodeType],
	fnNode Node[NodeType],
	bodyNode *BodyNode,
	ctx *parser.RevertStatementContext,
) Node[NodeType] {
	r.Src = SrcNode{
		Line:        int64(ctx.GetStart().GetLine()),
		Column:      int64(ctx.GetStart().GetColumn()),
		Start:       int64(ctx.GetStart().GetStart()),
		End:         int64(ctx.GetStop().GetStop()),
		Length:      int64(ctx.GetStop().GetStop() - ctx.GetStart().GetStart() + 1),
		ParentIndex: fnNode.GetId(),
	}

	expression := NewExpression(r.ASTBuilder)

	if ctx.CallArgumentList() != nil {
		for _, expressionCtx := range ctx.CallArgumentList().AllExpression() {
			r.Arguments = append(
				r.Arguments,
				expression.Parse(
					unit, contractNode, fnNode,
					bodyNode, nil, r, r.GetId(), expressionCtx,
				),
			)
		}
	}

	r.Expression = expression.Parse(unit, contractNode, fnNode, bodyNode, nil, r, r.GetId(), ctx.Expression())
	return r
}
