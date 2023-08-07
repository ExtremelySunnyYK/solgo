package ast

import (
	"fmt"

	ast_pb "github.com/txpull/protos/dist/go/ast"
	"github.com/txpull/solgo/parser"
)

type Expression struct {
	*ASTBuilder
}

func NewExpression(b *ASTBuilder) *Expression {
	return &Expression{
		ASTBuilder: b,
	}
}

func (e *Expression) Parse(
	unit *SourceUnit[Node[ast_pb.SourceUnit]],
	contractNode Node[NodeType],
	fnNode Node[NodeType],
	bodyNode *BodyNode,
	vDecar *VariableDeclaration,
	exprNode Node[NodeType],
	ctx parser.IExpressionContext,
) Node[NodeType] {
	switch ctxType := ctx.(type) {
	case *parser.AddSubOperationContext:
		binaryExp := NewBinaryOperationExpression(e.ASTBuilder)
		return binaryExp.ParseAddSub(unit, contractNode, fnNode, bodyNode, vDecar, exprNode, ctxType)
	case *parser.OrderComparisonContext:
		binaryExp := NewBinaryOperationExpression(e.ASTBuilder)
		return binaryExp.ParseOrderComparison(unit, contractNode, fnNode, bodyNode, vDecar, exprNode, ctxType)
	case *parser.MulDivModOperationContext:
		binaryExp := NewBinaryOperationExpression(e.ASTBuilder)
		return binaryExp.ParseMulDivMod(unit, contractNode, fnNode, bodyNode, vDecar, exprNode, ctxType)
	case *parser.EqualityComparisonContext:
		binaryExp := NewBinaryOperationExpression(e.ASTBuilder)
		return binaryExp.ParseEqualityComparison(unit, contractNode, fnNode, bodyNode, vDecar, exprNode, ctxType)
	case *parser.OrOperationContext:
		binaryExp := NewBinaryOperationExpression(e.ASTBuilder)
		return binaryExp.ParseOr(unit, contractNode, fnNode, bodyNode, vDecar, exprNode, ctxType)
	case *parser.AssignmentContext:
		assignment := NewAssignment(e.ASTBuilder)
		return assignment.Parse(unit, contractNode, fnNode, bodyNode, vDecar, exprNode, ctxType)
	case *parser.FunctionCallContext:
		statementNode := NewFunctionCall(e.ASTBuilder)
		return statementNode.Parse(unit, contractNode, fnNode, bodyNode, vDecar, exprNode, ctxType)
	case *parser.MemberAccessContext:
		memberAccess := NewMemberAccessExpression(e.ASTBuilder)
		return memberAccess.Parse(unit, contractNode, fnNode, bodyNode, vDecar, exprNode, ctxType)
	case *parser.PrimaryExpressionContext:
		primaryExp := NewPrimaryExpression(e.ASTBuilder)
		return primaryExp.Parse(unit, contractNode, fnNode, bodyNode, vDecar, exprNode, ctxType)
	case *parser.IndexAccessContext:
		indexAccess := NewIndexAccess(e.ASTBuilder)
		return indexAccess.Parse(unit, contractNode, fnNode, bodyNode, vDecar, exprNode, ctxType)
	case *parser.MetaTypeContext:
		metaType := NewMetaTypeExpression(e.ASTBuilder)
		return metaType.Parse(unit, contractNode, fnNode, bodyNode, vDecar, exprNode, ctxType)
	case *parser.PayableConversionContext:
		payableConversion := NewPayableConversionExpression(e.ASTBuilder)
		return payableConversion.Parse(unit, contractNode, fnNode, bodyNode, vDecar, exprNode, ctxType)
	case *parser.UnarySuffixOperationContext:
		unarySuffixOperation := NewUnarySuffixExpression(e.ASTBuilder)
		return unarySuffixOperation.Parse(unit, contractNode, fnNode, bodyNode, vDecar, exprNode, ctxType)
	case *parser.UnaryPrefixOperationContext:
		unaryPrefixOperation := NewUnaryPrefixExpression(e.ASTBuilder)
		return unaryPrefixOperation.Parse(unit, contractNode, fnNode, bodyNode, vDecar, exprNode, ctxType)
	case *parser.NewExprContext:
		newExpr := NewExprExpression(e.ASTBuilder)
		return newExpr.Parse(unit, contractNode, fnNode, bodyNode, vDecar, exprNode, ctxType)
	case *parser.TupleContext:
		tupleExpr := NewTupleExpression(e.ASTBuilder)
		return tupleExpr.Parse(unit, contractNode, fnNode, bodyNode, vDecar, exprNode, ctxType)
	case *parser.FunctionCallOptionsContext:
		statementNode := NewFunctionCallOption(e.ASTBuilder)
		return statementNode.Parse(unit, contractNode, fnNode, bodyNode, vDecar, exprNode, ctxType)
	case *parser.IndexRangeAccessContext:
		indexRangeAccess := NewIndexRangeAccessExpression(e.ASTBuilder)
		return indexRangeAccess.Parse(unit, contractNode, fnNode, bodyNode, vDecar, exprNode, ctxType)
	case *parser.ExpOperationContext:
		expOperation := NewExprOperationExpression(e.ASTBuilder)
		return expOperation.Parse(unit, contractNode, fnNode, bodyNode, vDecar, exprNode, ctxType)
	default:
		panic(
			fmt.Sprintf(
				"Expression type not supported @ Expression.Parse: %T",
				ctx,
			),
		)
	}
}
