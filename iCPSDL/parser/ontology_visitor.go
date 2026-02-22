// Code generated from .\iCPSDL\parser\ontology.g4 by ANTLR 4.10.1. DO NOT EDIT.

package parser // ontology

import "github.com/antlr/antlr4/runtime/Go/antlr"

// A complete Visitor for a parse tree produced by ontologyParser.
type ontologyVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by ontologyParser#program.
	VisitProgram(ctx *ProgramContext) interface{}

	// Visit a parse tree produced by ontologyParser#exit.
	VisitExit(ctx *ExitContext) interface{}

	// Visit a parse tree produced by ontologyParser#expression.
	VisitExpression(ctx *ExpressionContext) interface{}

	// Visit a parse tree produced by ontologyParser#load.
	VisitLoad(ctx *LoadContext) interface{}

	// Visit a parse tree produced by ontologyParser#assignment.
	VisitAssignment(ctx *AssignmentContext) interface{}

	// Visit a parse tree produced by ontologyParser#remove.
	VisitRemove(ctx *RemoveContext) interface{}

	// Visit a parse tree produced by ontologyParser#show.
	VisitShow(ctx *ShowContext) interface{}

	// Visit a parse tree produced by ontologyParser#mermaid.
	VisitMermaid(ctx *MermaidContext) interface{}

	// Visit a parse tree produced by ontologyParser#translate.
	VisitTranslate(ctx *TranslateContext) interface{}

	// Visit a parse tree produced by ontologyParser#traverse.
	VisitTraverse(ctx *TraverseContext) interface{}

	// Visit a parse tree produced by ontologyParser#configure.
	VisitConfigure(ctx *ConfigureContext) interface{}

	// Visit a parse tree produced by ontologyParser#compose.
	VisitCompose(ctx *ComposeContext) interface{}

	// Visit a parse tree produced by ontologyParser#project.
	VisitProject(ctx *ProjectContext) interface{}

	// Visit a parse tree produced by ontologyParser#info.
	VisitInfo(ctx *InfoContext) interface{}

	// Visit a parse tree produced by ontologyParser#domain.
	VisitDomain(ctx *DomainContext) interface{}

	// Visit a parse tree produced by ontologyParser#domain_declaration.
	VisitDomain_declaration(ctx *Domain_declarationContext) interface{}

	// Visit a parse tree produced by ontologyParser#property.
	VisitProperty(ctx *PropertyContext) interface{}

	// Visit a parse tree produced by ontologyParser#type.
	VisitType(ctx *TypeContext) interface{}

	// Visit a parse tree produced by ontologyParser#model.
	VisitModel(ctx *ModelContext) interface{}

	// Visit a parse tree produced by ontologyParser#class.
	VisitClass(ctx *ClassContext) interface{}

	// Visit a parse tree produced by ontologyParser#internal_edge.
	VisitInternal_edge(ctx *Internal_edgeContext) interface{}

	// Visit a parse tree produced by ontologyParser#translation.
	VisitTranslation(ctx *TranslationContext) interface{}

	// Visit a parse tree produced by ontologyParser#arg_connection.
	VisitArg_connection(ctx *Arg_connectionContext) interface{}

	// Visit a parse tree produced by ontologyParser#connection.
	VisitConnection(ctx *ConnectionContext) interface{}

	// Visit a parse tree produced by ontologyParser#process.
	VisitProcess(ctx *ProcessContext) interface{}

	// Visit a parse tree produced by ontologyParser#process_declaration.
	VisitProcess_declaration(ctx *Process_declarationContext) interface{}

	// Visit a parse tree produced by ontologyParser#device.
	VisitDevice(ctx *DeviceContext) interface{}

	// Visit a parse tree produced by ontologyParser#component.
	VisitComponent(ctx *ComponentContext) interface{}

	// Visit a parse tree produced by ontologyParser#connection_decl.
	VisitConnection_decl(ctx *Connection_declContext) interface{}

	// Visit a parse tree produced by ontologyParser#local_configuration.
	VisitLocal_configuration(ctx *Local_configurationContext) interface{}

	// Visit a parse tree produced by ontologyParser#global_configuration.
	VisitGlobal_configuration(ctx *Global_configurationContext) interface{}

	// Visit a parse tree produced by ontologyParser#local.
	VisitLocal(ctx *LocalContext) interface{}

	// Visit a parse tree produced by ontologyParser#send.
	VisitSend(ctx *SendContext) interface{}

	// Visit a parse tree produced by ontologyParser#receive.
	VisitReceive(ctx *ReceiveContext) interface{}

	// Visit a parse tree produced by ontologyParser#loop.
	VisitLoop(ctx *LoopContext) interface{}

	// Visit a parse tree produced by ontologyParser#label.
	VisitLabel(ctx *LabelContext) interface{}

	// Visit a parse tree produced by ontologyParser#select.
	VisitSelect(ctx *SelectContext) interface{}

	// Visit a parse tree produced by ontologyParser#branch.
	VisitBranch(ctx *BranchContext) interface{}

	// Visit a parse tree produced by ontologyParser#global.
	VisitGlobal(ctx *GlobalContext) interface{}

	// Visit a parse tree produced by ontologyParser#pass.
	VisitPass(ctx *PassContext) interface{}

	// Visit a parse tree produced by ontologyParser#global_loop.
	VisitGlobal_loop(ctx *Global_loopContext) interface{}

	// Visit a parse tree produced by ontologyParser#choice.
	VisitChoice(ctx *ChoiceContext) interface{}

	// Visit a parse tree produced by ontologyParser#repository.
	VisitRepository(ctx *RepositoryContext) interface{}

	// Visit a parse tree produced by ontologyParser#repository_decl.
	VisitRepository_decl(ctx *Repository_declContext) interface{}
}
