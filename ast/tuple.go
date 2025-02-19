package ast

import (
	"github.com/goccy/go-json"
	"strings"

	v3 "github.com/cncf/xds/go/xds/type/v3"
	ast_pb "github.com/unpackdev/protos/dist/go/ast"
	"github.com/unpackdev/solgo/parser"
)

// TupleExpression represents a tuple expression in Solidity.
type TupleExpression struct {
	*ASTBuilder                            // Embedding the ASTBuilder to provide common functionality
	Id                    int64            `json:"id"`                               // Unique identifier for the tuple expression
	NodeType              ast_pb.NodeType  `json:"node_type"`                        // Type of the node (TUPLE_EXPRESSION for a tuple expression)
	Src                   SrcNode          `json:"src"`                              // Source information about the tuple expression
	Constant              bool             `json:"is_constant"`                      // Whether the tuple expression is constant
	Pure                  bool             `json:"is_pure"`                          // Whether the tuple expression is pure
	Components            []Node[NodeType] `json:"components"`                       // Components of the tuple expression
	ReferencedDeclaration int64            `json:"referenced_declaration,omitempty"` // Referenced declaration of the tuple expression
	TypeDescription       *TypeDescription `json:"type_description"`                 // Type description of the tuple expression
}

// NewTupleExpression creates a new TupleExpression instance.
func NewTupleExpression(b *ASTBuilder) *TupleExpression {
	return &TupleExpression{
		ASTBuilder: b,
		Id:         b.GetNextID(),
		NodeType:   ast_pb.NodeType_TUPLE_EXPRESSION,
	}
}

// SetReferenceDescriptor sets the reference descriptions of the TupleExpression node.
func (t *TupleExpression) SetReferenceDescriptor(refId int64, refDesc *TypeDescription) bool {
	t.ReferencedDeclaration = refId
	t.TypeDescription = refDesc
	return false
}

// GetId returns the unique identifier of the tuple expression.
func (t *TupleExpression) GetId() int64 {
	return t.Id
}

// GetType returns the type of the node, which is 'TUPLE_EXPRESSION' for a tuple expression.
func (t *TupleExpression) GetType() ast_pb.NodeType {
	return t.NodeType
}

// GetSrc returns the source information about the tuple expression.
func (t *TupleExpression) GetSrc() SrcNode {
	return t.Src
}

// GetComponents returns the components of the tuple expression.
func (t *TupleExpression) GetComponents() []Node[NodeType] {
	return t.Components
}

// GetNodes returns the components of the tuple expression.
func (t *TupleExpression) GetNodes() []Node[NodeType] {
	return t.Components
}

// GetTypeDescription returns the type description of the tuple expression.
func (t *TupleExpression) GetTypeDescription() *TypeDescription {
	return t.TypeDescription
}

// IsConstant returns whether the tuple expression is constant.
func (t *TupleExpression) IsConstant() bool {
	return t.Constant
}

// IsPure returns whether the tuple expression is pure.
func (t *TupleExpression) IsPure() bool {
	return t.Pure
}

// GetReferencedDeclaration returns the referenced declaration of the tuple expression.
func (t *TupleExpression) GetReferencedDeclaration() int64 {
	return t.ReferencedDeclaration
}

// MarshalJSON marshals the TupleExpression node into a JSON byte slice.
func (t *TupleExpression) UnmarshalJSON(data []byte) error {
	var tempMap map[string]json.RawMessage
	if err := json.Unmarshal(data, &tempMap); err != nil {
		return err
	}

	if id, ok := tempMap["id"]; ok {
		if err := json.Unmarshal(id, &t.Id); err != nil {
			return err
		}
	}

	if nodeType, ok := tempMap["node_type"]; ok {
		if err := json.Unmarshal(nodeType, &t.NodeType); err != nil {
			return err
		}
	}

	if src, ok := tempMap["src"]; ok {
		if err := json.Unmarshal(src, &t.Src); err != nil {
			return err
		}
	}

	if constant, ok := tempMap["is_constant"]; ok {
		if err := json.Unmarshal(constant, &t.Constant); err != nil {
			return err
		}
	}

	if pure, ok := tempMap["is_pure"]; ok {
		if err := json.Unmarshal(pure, &t.Pure); err != nil {
			return err
		}
	}

	if referencedDeclaration, ok := tempMap["referenced_declaration"]; ok {
		if err := json.Unmarshal(referencedDeclaration, &t.ReferencedDeclaration); err != nil {
			return err
		}
	}

	if typeDescription, ok := tempMap["type_description"]; ok {
		if err := json.Unmarshal(typeDescription, &t.TypeDescription); err != nil {
			return err
		}
	}

	if components, ok := tempMap["components"]; ok {
		var nodes []json.RawMessage
		if err := json.Unmarshal(components, &nodes); err != nil {
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
			t.Components = append(t.Components, node)
		}
	}

	return nil
}

// ToProto returns the protobuf representation of the tuple expression.
func (t *TupleExpression) ToProto() NodeType {
	proto := ast_pb.Tuple{
		Id:                    t.GetId(),
		NodeType:              t.GetType(),
		Src:                   t.GetSrc().ToProto(),
		IsConstant:            t.IsConstant(),
		IsPure:                t.IsPure(),
		ReferencedDeclaration: t.GetReferencedDeclaration(),
		TypeDescription:       t.GetTypeDescription().ToProto(),
	}

	for _, component := range t.GetComponents() {
		proto.Components = append(proto.Components, component.ToProto().(*v3.TypedStruct))
	}

	return NewTypedStruct(&proto, "Tuple")
}

// Parse parses a tuple expression from the provided parser.TupleContext and returns the corresponding TupleExpression.
func (t *TupleExpression) Parse(
	unit *SourceUnit[Node[ast_pb.SourceUnit]],
	contractNode Node[NodeType],
	fnNode Node[NodeType],
	bodyNode *BodyNode,
	vDeclar *VariableDeclaration,
	exprNode Node[NodeType],
	ctx *parser.TupleContext,
) Node[NodeType] {
	t.Src = SrcNode{
		Line:   int64(ctx.GetStart().GetLine()),
		Column: int64(ctx.GetStart().GetColumn()),
		Start:  int64(ctx.GetStart().GetStart()),
		End:    int64(ctx.GetStop().GetStop()),
		Length: int64(ctx.GetStop().GetStop() - ctx.GetStart().GetStart() + 1),
		ParentIndex: func() int64 {
			if exprNode != nil {
				return exprNode.GetId()
			}

			if fnNode != nil {
				return fnNode.GetId()
			}

			return bodyNode.GetId()
		}(),
	}

	expression := NewExpression(t.ASTBuilder)
	for _, tupleCtx := range ctx.TupleExpression().AllExpression() {
		expr := expression.Parse(unit, contractNode, fnNode, bodyNode, vDeclar, t, t.GetId(), tupleCtx)
		t.Components = append(
			t.Components,
			expr,
		)
		// A bit of a hack as we have interfaces but it works...
		switch exprCtx := expr.(type) {
		case *PrimaryExpression:
			if exprCtx.IsPure() {
				t.Pure = true
				break
			}
		}
	}

	t.TypeDescription = t.buildTypeDescription()
	return t
}

// buildTypeDescription constructs the type description of the tuple expression.
func (t *TupleExpression) buildTypeDescription() *TypeDescription {
	typeString := "tuple("
	typeIdentifier := "t_tuple_"
	typeStrings := make([]string, 0)
	typeIdentifiers := make([]string, 0)

	for _, component := range t.GetComponents() {
		td := component.GetTypeDescription()
		if td == nil {
			typeStrings = append(typeStrings, "unknown")
			typeIdentifiers = append(typeIdentifiers, "$_t_unknown")
			continue
		}
		typeStrings = append(typeStrings, td.TypeString)
		typeIdentifiers = append(typeIdentifiers, "$_"+td.TypeIdentifier)
	}

	typeString += strings.Join(typeStrings, ",") + ")"
	typeIdentifier += strings.Join(typeIdentifiers, "_")
	typeIdentifier += "$"

	return &TypeDescription{
		TypeString:     typeString,
		TypeIdentifier: typeIdentifier,
	}
}
