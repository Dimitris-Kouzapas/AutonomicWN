// Code generated from .\iCPSDL\parser\ontology.g4 by ANTLR 4.10.1. DO NOT EDIT.

package parser // ontology

import "github.com/antlr/antlr4/runtime/Go/antlr"

type BaseontologyVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseontologyVisitor) VisitProgram(ctx *ProgramContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseontologyVisitor) VisitExit(ctx *ExitContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseontologyVisitor) VisitExpression(ctx *ExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseontologyVisitor) VisitLoad(ctx *LoadContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseontologyVisitor) VisitAssignment(ctx *AssignmentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseontologyVisitor) VisitRemove(ctx *RemoveContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseontologyVisitor) VisitShow(ctx *ShowContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseontologyVisitor) VisitMermaid(ctx *MermaidContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseontologyVisitor) VisitTranslate(ctx *TranslateContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseontologyVisitor) VisitTraverse(ctx *TraverseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseontologyVisitor) VisitConfigure(ctx *ConfigureContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseontologyVisitor) VisitCompose(ctx *ComposeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseontologyVisitor) VisitProject(ctx *ProjectContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseontologyVisitor) VisitInfo(ctx *InfoContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseontologyVisitor) VisitDomain(ctx *DomainContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseontologyVisitor) VisitDomain_declaration(ctx *Domain_declarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseontologyVisitor) VisitProperty(ctx *PropertyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseontologyVisitor) VisitType(ctx *TypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseontologyVisitor) VisitModel(ctx *ModelContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseontologyVisitor) VisitClass(ctx *ClassContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseontologyVisitor) VisitInternal_edge(ctx *Internal_edgeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseontologyVisitor) VisitTranslation(ctx *TranslationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseontologyVisitor) VisitArg_connection(ctx *Arg_connectionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseontologyVisitor) VisitConnection(ctx *ConnectionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseontologyVisitor) VisitProcess(ctx *ProcessContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseontologyVisitor) VisitProcess_declaration(ctx *Process_declarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseontologyVisitor) VisitDevice(ctx *DeviceContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseontologyVisitor) VisitComponent(ctx *ComponentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseontologyVisitor) VisitConnection_decl(ctx *Connection_declContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseontologyVisitor) VisitLocal_configuration(ctx *Local_configurationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseontologyVisitor) VisitGlobal_configuration(ctx *Global_configurationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseontologyVisitor) VisitLocal(ctx *LocalContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseontologyVisitor) VisitSend(ctx *SendContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseontologyVisitor) VisitReceive(ctx *ReceiveContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseontologyVisitor) VisitLoop(ctx *LoopContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseontologyVisitor) VisitLabel(ctx *LabelContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseontologyVisitor) VisitSelect(ctx *SelectContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseontologyVisitor) VisitBranch(ctx *BranchContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseontologyVisitor) VisitGlobal(ctx *GlobalContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseontologyVisitor) VisitPass(ctx *PassContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseontologyVisitor) VisitGlobal_loop(ctx *Global_loopContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseontologyVisitor) VisitChoice(ctx *ChoiceContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseontologyVisitor) VisitRepository(ctx *RepositoryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseontologyVisitor) VisitRepository_decl(ctx *Repository_declContext) interface{} {
	return v.VisitChildren(ctx)
}
