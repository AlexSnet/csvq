package query

import (
	"fmt"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/ternary"
)

type Procedure struct {
	Filter    *Filter
	ReturnVal parser.Primary
}

func NewProcedure() *Procedure {
	return &Procedure{
		Filter: NewEmptyFilter(),
	}
}

func (proc *Procedure) NewChildProcedure() *Procedure {
	return &Procedure{
		Filter: proc.Filter.CreateChildScope(),
	}
}

func (proc *Procedure) ExecuteChild(statements []parser.Statement) (StatementFlow, error) {
	child := proc.NewChildProcedure()
	f, err := child.Execute(statements)
	return f, err
}

func (proc *Procedure) Execute(statements []parser.Statement) (StatementFlow, error) {
	flow := TERMINATE

	for _, stmt := range statements {
		f, err := proc.ExecuteStatement(stmt)
		if err != nil {
			return f, err
		}
		if f != TERMINATE {
			flow = f
			break
		}
	}
	return flow, nil
}

func (proc *Procedure) ExecuteStatement(stmt parser.Statement) (StatementFlow, error) {
	flow := TERMINATE

	var err error

	var results []Result
	var view *View
	var views []*View
	var printstr string

	switch stmt.(type) {
	case parser.SetFlag:
		err = SetFlag(stmt.(parser.SetFlag))
	case parser.VariableDeclaration:
		err = proc.Filter.VariablesList.Declare(stmt.(parser.VariableDeclaration), proc.Filter)
	case parser.VariableSubstitution:
		_, err = proc.Filter.Evaluate(stmt.(parser.QueryExpression))
	case parser.DisposeVariable:
		err = proc.Filter.VariablesList.Dispose(stmt.(parser.DisposeVariable).Variable)
	case parser.CursorDeclaration:
		err = proc.Filter.CursorsList.Declare(stmt.(parser.CursorDeclaration))
	case parser.OpenCursor:
		err = proc.Filter.CursorsList.Open(stmt.(parser.OpenCursor).Cursor, proc.Filter)
	case parser.CloseCursor:
		err = proc.Filter.CursorsList.Close(stmt.(parser.CloseCursor).Cursor)
	case parser.DisposeCursor:
		err = proc.Filter.CursorsList.Dispose(stmt.(parser.DisposeCursor).Cursor)
	case parser.FetchCursor:
		fetch := stmt.(parser.FetchCursor)
		_, err = FetchCursor(fetch.Cursor, fetch.Position, fetch.Variables, proc.Filter)
	case parser.TableDeclaration:
		err = DeclareTable(stmt.(parser.TableDeclaration), proc.Filter)
	case parser.DisposeTable:
		err = proc.Filter.TempViewsList.Dispose(stmt.(parser.DisposeTable).Table)
	case parser.FunctionDeclaration:
		err = proc.Filter.FunctionsList.Declare(stmt.(parser.FunctionDeclaration))
	case parser.AggregateDeclaration:
		err = proc.Filter.FunctionsList.DeclareAggregate(stmt.(parser.AggregateDeclaration))
	case parser.SelectQuery:
		if view, err = Select(stmt.(parser.SelectQuery), proc.Filter); err == nil {
			flags := cmd.GetFlags()
			var viewstr string
			var lineBreak = cmd.LF
			if 0 < len(flags.OutFile) {
				lineBreak = flags.LineBreak
			}
			viewstr, err = EncodeView(view, flags.Format, flags.WriteDelimiter, flags.WithoutHeader, flags.WriteEncoding, lineBreak)
			if err == nil {
				if 0 < len(flags.OutFile) {
					AddSelectLog(viewstr)
				} else {
					Log(viewstr)
				}
			}
		}
	case parser.InsertQuery:
		if view, err = Insert(stmt.(parser.InsertQuery), proc.Filter); err == nil {
			results = []Result{
				{
					Type:          INSERT,
					FileInfo:      view.FileInfo,
					OperatedCount: view.OperatedRecords,
				},
			}
			Log(fmt.Sprintf("%s inserted on %q.", FormatCount(view.OperatedRecords, "record"), view.FileInfo.Path))

			view.OperatedRecords = 0
		}
	case parser.UpdateQuery:
		if views, err = Update(stmt.(parser.UpdateQuery), proc.Filter); err == nil {
			results = make([]Result, len(views))
			for i, v := range views {
				results[i] = Result{
					Type:          UPDATE,
					FileInfo:      v.FileInfo,
					OperatedCount: v.OperatedRecords,
				}
				Log(fmt.Sprintf("%s updated on %q.", FormatCount(v.OperatedRecords, "record"), v.FileInfo.Path))

				v.OperatedRecords = 0
			}
		}
	case parser.DeleteQuery:
		if views, err = Delete(stmt.(parser.DeleteQuery), proc.Filter); err == nil {
			results = make([]Result, len(views))
			for i, v := range views {
				results[i] = Result{
					Type:          DELETE,
					FileInfo:      v.FileInfo,
					OperatedCount: v.OperatedRecords,
				}
				Log(fmt.Sprintf("%s deleted on %q.", FormatCount(v.OperatedRecords, "record"), v.FileInfo.Path))

				v.OperatedRecords = 0
			}
		}
	case parser.CreateTable:
		if view, err = CreateTable(stmt.(parser.CreateTable), proc.Filter); err == nil {
			results = []Result{
				{
					Type:     CREATE_TABLE,
					FileInfo: view.FileInfo,
				},
			}
			Log(fmt.Sprintf("file %q is created.", view.FileInfo.Path))

			view.OperatedRecords = 0
		}
	case parser.AddColumns:
		if view, err = AddColumns(stmt.(parser.AddColumns), proc.Filter); err == nil {
			results = []Result{
				{
					Type:          ADD_COLUMNS,
					FileInfo:      view.FileInfo,
					OperatedCount: view.OperatedFields,
				},
			}
			Log(fmt.Sprintf("%s added on %q.", FormatCount(view.OperatedFields, "field"), view.FileInfo.Path))

			view.OperatedRecords = 0
		}
	case parser.DropColumns:
		if view, err = DropColumns(stmt.(parser.DropColumns), proc.Filter); err == nil {
			results = []Result{
				{
					Type:          DROP_COLUMNS,
					FileInfo:      view.FileInfo,
					OperatedCount: view.OperatedFields,
				},
			}
			Log(fmt.Sprintf("%s dropped on %q.", FormatCount(view.OperatedFields, "field"), view.FileInfo.Path))

			view.OperatedRecords = 0
		}
	case parser.RenameColumn:
		if view, err = RenameColumn(stmt.(parser.RenameColumn), proc.Filter); err == nil {
			results = []Result{
				{
					Type:          RENAME_COLUMN,
					FileInfo:      view.FileInfo,
					OperatedCount: view.OperatedFields,
				},
			}
			Log(fmt.Sprintf("%s renamed on %q.", FormatCount(view.OperatedFields, "field"), view.FileInfo.Path))

			view.OperatedRecords = 0
		}
	case parser.TransactionControl:
		switch stmt.(parser.TransactionControl).Token {
		case parser.COMMIT:
			err = proc.Commit(stmt.(parser.Expression))
		case parser.ROLLBACK:
			proc.Rollback()
		}
	case parser.FlowControl:
		switch stmt.(parser.FlowControl).Token {
		case parser.CONTINUE:
			flow = CONTINUE
		case parser.BREAK:
			flow = BREAK
		case parser.EXIT:
			flow = EXIT
		}
	case parser.Return:
		var ret parser.Primary
		if ret, err = proc.Filter.Evaluate(stmt.(parser.Return).Value); err == nil {
			proc.ReturnVal = ret
			flow = RETURN
		}
	case parser.If:
		flow, err = proc.IfStmt(stmt.(parser.If))
	case parser.Case:
		flow, err = proc.Case(stmt.(parser.Case))
	case parser.While:
		flow, err = proc.While(stmt.(parser.While))
	case parser.WhileInCursor:
		flow, err = proc.WhileInCursor(stmt.(parser.WhileInCursor))
	case parser.Print:
		if printstr, err = Print(stmt.(parser.Print), proc.Filter); err == nil {
			Log(printstr)
		}
	case parser.Function:
		_, err = proc.Filter.Evaluate(stmt.(parser.Function))
	case parser.Printf:
		if printstr, err = Printf(stmt.(parser.Printf), proc.Filter); err == nil {
			Log(printstr)
		}
	case parser.Source:
		var externalStatements []parser.Statement
		source := stmt.(parser.Source)
		if externalStatements, err = Source(source, proc.Filter); err == nil {
			flow, err = proc.Execute(externalStatements)
		}
	case parser.Trigger:
		trigger := stmt.(parser.Trigger)
		switch trigger.Token {
		case parser.ERROR:
			var message string
			if trigger.Message != nil {
				if pt, ok := trigger.Message.(parser.PrimitiveType); ok && trigger.Code == nil && pt.IsInteger() {
					trigger.Code = pt.Value
				} else {
					var p parser.Primary
					if p, err = proc.Filter.Evaluate(trigger.Message); err == nil {
						if s := parser.PrimaryToString(p); !parser.IsNull(s) {
							message = s.(parser.String).Value()
						}
					}
				}
			}
			if err == nil {
				err = NewUserTriggeredError(trigger, message)
			}
		}
	}

	if results != nil {
		Results = append(Results, results...)
	}

	if err != nil {
		flow = ERROR
	}
	return flow, err
}

func (proc *Procedure) IfStmt(stmt parser.If) (StatementFlow, error) {
	stmts := make([]parser.ElseIf, len(stmt.ElseIf)+1)
	stmts[0] = parser.ElseIf{
		Condition:  stmt.Condition,
		Statements: stmt.Statements,
	}
	for i, v := range stmt.ElseIf {
		stmts[i+1] = v.(parser.ElseIf)
	}

	for _, v := range stmts {
		p, err := proc.Filter.Evaluate(v.Condition)
		if err != nil {
			return ERROR, err
		}
		if p.Ternary() == ternary.TRUE {
			return proc.ExecuteChild(v.Statements)
		}
	}

	if stmt.Else != nil {
		return proc.ExecuteChild(stmt.Else.(parser.Else).Statements)
	}
	return TERMINATE, nil
}

func (proc *Procedure) Case(stmt parser.Case) (StatementFlow, error) {
	var value parser.Primary
	var err error
	if stmt.Value != nil {
		value, err = proc.Filter.Evaluate(stmt.Value)
		if err != nil {
			return ERROR, err
		}
	}

	for _, v := range stmt.When {
		when := v.(parser.CaseWhen)
		var t ternary.Value

		cond, err := proc.Filter.Evaluate(when.Condition)
		if err != nil {
			return ERROR, err
		}

		if value == nil {
			t = cond.Ternary()
		} else {
			t = EqualTo(value, cond)
		}

		if t == ternary.TRUE {
			return proc.ExecuteChild(when.Statements)
		}
	}

	if stmt.Else == nil {
		return TERMINATE, nil
	}
	return proc.ExecuteChild(stmt.Else.(parser.CaseElse).Statements)
}

func (proc *Procedure) While(stmt parser.While) (StatementFlow, error) {
	for {
		p, err := proc.Filter.Evaluate(stmt.Condition)
		if err != nil {
			return ERROR, err
		}
		if p.Ternary() != ternary.TRUE {
			break
		}
		f, err := proc.ExecuteChild(stmt.Statements)
		if err != nil {
			return ERROR, err
		}

		if f == BREAK {
			return TERMINATE, nil
		}
		if f == EXIT {
			return EXIT, nil
		}
	}
	return TERMINATE, nil
}

func (proc *Procedure) WhileInCursor(stmt parser.WhileInCursor) (StatementFlow, error) {
	for {
		success, err := FetchCursor(stmt.Cursor, nil, stmt.Variables, proc.Filter)
		if err != nil {
			return ERROR, err
		}
		if !success {
			break
		}

		f, err := proc.ExecuteChild(stmt.Statements)
		if err != nil {
			return ERROR, err
		}

		if f == BREAK {
			return TERMINATE, nil
		}
		if f == EXIT {
			return EXIT, nil
		}
	}

	return TERMINATE, nil
}

func (proc *Procedure) Commit(expr parser.Expression) error {
	var createFiles = map[string]*FileInfo{}
	var updateFiles = map[string]*FileInfo{}

	for _, result := range Results {
		if result.FileInfo != nil {
			switch result.Type {
			case CREATE_TABLE:
				createFiles[result.FileInfo.Path] = result.FileInfo
			default:
				if !result.FileInfo.IsTemporary && 0 < result.OperatedCount {
					if _, ok := createFiles[result.FileInfo.Path]; !ok {
						if _, ok := updateFiles[result.FileInfo.Path]; !ok {
							updateFiles[result.FileInfo.Path] = result.FileInfo
						}
					}
				}
			}
		}
	}

	var modified bool

	if 0 < len(createFiles) {
		for filename, fileinfo := range createFiles {
			view, _ := ViewCache.Get(parser.Identifier{Literal: filename})
			viewstr, err := EncodeView(view, cmd.CSV, fileinfo.Delimiter, false, fileinfo.Encoding, fileinfo.LineBreak)
			if err != nil {
				return err
			}

			if err = cmd.CreateFile(filename, viewstr); err != nil {
				if expr == nil {
					return NewAutoCommitError(err.Error())
				}
				return NewWriteFileError(expr, err.Error())
			}
			Log(fmt.Sprintf("Commit: file %q is created.", filename))
			if !modified {
				modified = true
			}
		}
	}

	if 0 < len(updateFiles) {
		for filename, fileinfo := range updateFiles {
			view, _ := ViewCache.Get(parser.Identifier{Literal: filename})
			viewstr, err := EncodeView(view, cmd.CSV, fileinfo.Delimiter, fileinfo.NoHeader, fileinfo.Encoding, fileinfo.LineBreak)
			if err != nil {
				return err
			}

			if err = cmd.UpdateFile(filename, viewstr); err != nil {
				if expr == nil {
					return NewAutoCommitError(err.Error())
				}
				return NewWriteFileError(expr, err.Error())
			}
			Log(fmt.Sprintf("Commit: file %q is updated.", filename))
			if !modified {
				modified = true
			}
		}
	}

	Results = []Result{}
	ViewCache.Clear()

	return nil
}

func (proc *Procedure) Rollback() {
	Results = []Result{}
	ViewCache.Clear()
	proc.Filter.TempViewsList.Rollback()

	Log("Rolled back.")
	return
}
