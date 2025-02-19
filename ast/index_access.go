package ast

import (
	"fmt"
	"github.com/goccy/go-json"
	"regexp"
	"strings"

	v3 "github.com/cncf/xds/go/xds/type/v3"
	ast_pb "github.com/unpackdev/protos/dist/go/ast"
	"github.com/unpackdev/solgo/parser"
)

// IndexAccess represents an index access expression in the AST.
type IndexAccess struct {
	*ASTBuilder

	Id                    int64              `json:"id"`                               // Unique identifier for the IndexAccess node.
	NodeType              ast_pb.NodeType    `json:"node_type"`                        // Type of the AST node.
	Src                   SrcNode            `json:"src"`                              // Source location information.
	IndexExpression       Node[NodeType]     `json:"index_expression"`                 // Index expression.
	BaseExpression        Node[NodeType]     `json:"base_expression"`                  // Base expression.
	TypeDescriptions      []*TypeDescription `json:"type_descriptions"`                // Type descriptions.
	ReferencedDeclaration int64              `json:"referenced_declaration,omitempty"` // Referenced declaration.
	TypeDescription       *TypeDescription   `json:"type_description"`                 // Type description.
}

// NewIndexAccess creates a new IndexAccess node with a given ASTBuilder.
func NewIndexAccess(b *ASTBuilder) *IndexAccess {
	return &IndexAccess{
		ASTBuilder: b,
		Id:         b.GetNextID(),
		NodeType:   ast_pb.NodeType_INDEX_ACCESS,
	}
}

// SetReferenceDescriptor sets the reference descriptions of the IndexAccess node.
// Here we are going to just do some magic stuff in order to figure out descriptions across the board...
func (i *IndexAccess) SetReferenceDescriptor(refId int64, refDesc *TypeDescription) bool {
	// It is usually only index expression that is affected, so for now fixing that one...
	if i.IndexExpression != nil && i.IndexExpression.GetTypeDescription() != nil {
		i.TypeDescriptions[0] = i.IndexExpression.GetTypeDescription()
		i.TypeDescription = i.buildTypeDescription()
		return true
	}
	return true
}

// GetName returns the name of the IndexAccess node.
func (i *IndexAccess) GetName() string {
	return fmt.Sprintf("index_access_%d", i.Id)
}

// GetId returns the ID of the IndexAccess node.
func (i *IndexAccess) GetId() int64 {
	return i.Id
}

// GetType returns the NodeType of the IndexAccess node.
func (i *IndexAccess) GetType() ast_pb.NodeType {
	return i.NodeType
}

// GetSrc returns the SrcNode of the IndexAccess node.
func (i *IndexAccess) GetSrc() SrcNode {
	return i.Src
}

// GetIndexExpression returns the index expression.
func (i *IndexAccess) GetIndexExpression() Node[NodeType] {
	return i.IndexExpression
}

// GetBaseExpression returns the base expression.
func (i *IndexAccess) GetBaseExpression() Node[NodeType] {
	return i.BaseExpression
}

// GetTypeDescription returns the type description.
func (i *IndexAccess) GetTypeDescription() *TypeDescription {
	return i.TypeDescription
}

// GetTypeDescriptions returns the list of type descriptions.
func (i *IndexAccess) GetTypeDescriptions() []*TypeDescription {
	return i.TypeDescriptions
}

// GetNodes returns the child nodes of the IndexAccess node.
func (i *IndexAccess) GetNodes() []Node[NodeType] {
	toReturn := []Node[NodeType]{i.IndexExpression}
	if i.BaseExpression != nil {
		toReturn = append(toReturn, i.BaseExpression)
	}
	return toReturn
}

// GetReferencedDeclaration returns the referenced declaration.
func (i *IndexAccess) GetReferencedDeclaration() int64 {
	return i.ReferencedDeclaration
}

// UnmarshalJSON sets the IndexAccess node data from its JSON representation.
func (i *IndexAccess) UnmarshalJSON(data []byte) error {
	var tempMap map[string]json.RawMessage
	if err := json.Unmarshal(data, &tempMap); err != nil {
		return err
	}

	if id, ok := tempMap["id"]; ok {
		if err := json.Unmarshal(id, &i.Id); err != nil {
			return err
		}
	}

	if nodeType, ok := tempMap["node_type"]; ok {
		if err := json.Unmarshal(nodeType, &i.NodeType); err != nil {
			return err
		}
	}

	if src, ok := tempMap["src"]; ok {
		if err := json.Unmarshal(src, &i.Src); err != nil {
			return err
		}
	}

	if indexExpression, ok := tempMap["index_expression"]; ok {
		if err := json.Unmarshal(indexExpression, &i.IndexExpression); err != nil {
			var tempNodeMap map[string]json.RawMessage
			if err := json.Unmarshal(indexExpression, &tempNodeMap); err != nil {
				return err
			}

			var tempNodeType ast_pb.NodeType
			if err := json.Unmarshal(tempNodeMap["node_type"], &tempNodeType); err != nil {
				return err
			}

			node, err := unmarshalNode(indexExpression, tempNodeType)
			if err != nil {
				return err
			}
			i.IndexExpression = node
		}
	}

	if baseExpression, ok := tempMap["base_expression"]; ok {
		if err := json.Unmarshal(baseExpression, &i.BaseExpression); err != nil {
			var tempNodeMap map[string]json.RawMessage
			if err := json.Unmarshal(baseExpression, &tempNodeMap); err != nil {
				return err
			}

			var tempNodeType ast_pb.NodeType
			if err := json.Unmarshal(tempNodeMap["node_type"], &tempNodeType); err != nil {
				return err
			}

			node, err := unmarshalNode(baseExpression, tempNodeType)
			if err != nil {
				return err
			}
			i.BaseExpression = node
		}
	}

	if typeDescriptions, ok := tempMap["type_descriptions"]; ok {
		if err := json.Unmarshal(typeDescriptions, &i.TypeDescriptions); err != nil {
			return err
		}
	}

	if referencedDeclaration, ok := tempMap["referenced_declaration"]; ok {
		if err := json.Unmarshal(referencedDeclaration, &i.ReferencedDeclaration); err != nil {
			return err
		}
	}

	if typeDescription, ok := tempMap["type_description"]; ok {
		if err := json.Unmarshal(typeDescription, &i.TypeDescription); err != nil {
			return err
		}
	}

	return nil
}

// ToProto returns a protobuf representation of the IndexAccess node.
func (i *IndexAccess) ToProto() NodeType {
	proto := ast_pb.IndexAccess{
		Id:                    i.GetId(),
		NodeType:              i.GetType(),
		Src:                   i.Src.ToProto(),
		TypeDescriptions:      make([]*ast_pb.TypeDescription, 0),
		ReferencedDeclaration: i.GetReferencedDeclaration(),
		TypeDescription:       i.GetTypeDescription().ToProto(),
	}

	if i.GetIndexExpression() != nil {
		proto.IndexExpression = i.GetIndexExpression().ToProto().(*v3.TypedStruct)
	}

	if i.GetBaseExpression() != nil {
		proto.BaseExpression = i.GetBaseExpression().ToProto().(*v3.TypedStruct)
	}

	for _, td := range i.GetTypeDescriptions() {
		if td != nil {
			proto.TypeDescriptions = append(proto.TypeDescriptions, td.ToProto())
		}
	}

	return NewTypedStruct(&proto, "IndexAccess")
}

// Parse parses an index access context into the IndexAccess node.
func (i *IndexAccess) Parse(
	unit *SourceUnit[Node[ast_pb.SourceUnit]],
	contractNode Node[NodeType],
	fnNode Node[NodeType],
	bodyNode *BodyNode,
	vDeclar *VariableDeclaration,
	expNode Node[NodeType],
	ctx *parser.IndexAccessContext,
) Node[NodeType] {
	i.Src = SrcNode{
		Line:   int64(ctx.GetStart().GetLine()),
		Column: int64(ctx.GetStart().GetColumn()),
		Start:  int64(ctx.GetStart().GetStart()),
		End:    int64(ctx.GetStop().GetStop()),
		Length: int64(ctx.GetStop().GetStop() - ctx.GetStart().GetStart() + 1),
		ParentIndex: func() int64 {
			if vDeclar != nil {
				return vDeclar.GetId()
			}

			if expNode != nil {
				return expNode.GetId()
			}

			return bodyNode.GetId()
		}(),
	}

	expression := NewExpression(i.ASTBuilder)

	if ctx.Expression(0) != nil {
		i.BaseExpression = expression.Parse(
			unit, contractNode, fnNode, bodyNode, vDeclar, i, i.GetId(), ctx.Expression(0),
		)
		i.TypeDescriptions = append(i.TypeDescriptions, i.BaseExpression.GetTypeDescription())
	}

	if ctx.Expression(1) != nil {
		i.IndexExpression = expression.Parse(
			unit, contractNode, fnNode, bodyNode, vDeclar, i, i.GetId(), ctx.Expression(1),
		)

		i.TypeDescription = i.IndexExpression.GetTypeDescription()

		i.TypeDescriptions = []*TypeDescription{
			i.IndexExpression.GetTypeDescription(),
		}
	}

	if i.IndexExpression != nil && i.IndexExpression.GetTypeDescription() == nil || (i.BaseExpression != nil && i.BaseExpression.GetTypeDescription() == nil) {
		if refId, refTypeDescription := i.GetResolver().ResolveByNode(i, fmt.Sprintf("index_access_%d", i.Id)); refTypeDescription != nil {
			i.ReferencedDeclaration = refId
			i.TypeDescription = refTypeDescription
			i.TypeDescription = i.buildTypeDescription()
		}
	}

	i.TypeDescription = i.buildTypeDescription()
	return i
}

// buildTypeDescription creates a type description for the IndexAccess node.
func (i *IndexAccess) buildTypeDescription() *TypeDescription {
	typeString := "index["
	typeIdentifier := "t_[_["
	typeStrings := make([]string, 0)
	typeIdentifiers := make([]string, 0)

	for _, paramType := range i.GetTypeDescriptions() {
		// REMOVE-LATER: It's a fix because sometimes forward-path is not quite working at this stage...
		// For example, defining state variables at end of the contract instead of the top :explosion:
		if paramType == nil {
			typeStrings = append(typeStrings, "unknown")
			typeIdentifiers = append(typeIdentifiers, "$_t_unknown")
			continue
		}

		if strings.Contains(paramType.TypeString, "literal_string") {
			typeStrings = append(typeStrings, "string memory")
			typeIdentifiers = append(typeIdentifiers, "_"+paramType.TypeIdentifier)
			continue
		} else if strings.Contains(paramType.TypeString, "contract") {
			typeStrings = append(typeStrings, "address")
			typeIdentifiers = append(typeIdentifiers, "$_t_address")
			continue
		}

		typeStrings = append(typeStrings, paramType.TypeString)
		typeIdentifiers = append(typeIdentifiers, "$_"+paramType.TypeIdentifier)
	}

	typeString += strings.Join(typeStrings, ":") + "]"
	typeIdentifier += strings.Join(typeIdentifiers, "]$")

	if !strings.HasSuffix(typeIdentifier, "$") {
		typeIdentifier += "]$"
	}

	re := regexp.MustCompile(`\${2,}`)
	typeIdentifier = re.ReplaceAllString(typeIdentifier, "$")

	return &TypeDescription{
		TypeString:     typeString,
		TypeIdentifier: typeIdentifier,
	}
}
