// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.648
package home

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	expense "mercado/app/models/expense"
	l "mercado/app/views/layouts"
)

func Index(expenses []expense.Expense) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<script src=\"public/js/home.js\"></script>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Var2 := templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
			templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
			if !templ_7745c5c3_IsBuffer {
				templ_7745c5c3_Buffer = templ.GetBuffer()
				defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<form id=\"add-form\" hx-target=\"#expenses-list\" hx-post=\"/\" hx-swap=\"afterbegin\" x-data=\"\n        {\n            expenseValue: undefined,\n            adding: false,\n        }\n    \" x-init=\"handleFormInit\" class=\"sticky top-0 mx-auto max-w-sm rounded-md bg-white p-4 shadow-md\"><div class=\"mb-4\"><label for=\"value\" class=\"mb-2 block text-sm font-bold text-gray-700\">Valor</label> <input id=\"value\" name=\"value\" placeholder=\"É de quanto, é de quanto?\" type=\"text\" inputmode=\"decimal\" x-model=\"expenseValue\" :value=\"expenseValue\" @input=\"handleInputChange\" class=\"w-full rounded-md border px-3 py-2 focus:border-blue-500 focus:outline-none\"></div><div class=\"mb-4\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = GroceryStoreSelector("", true).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><button id=\"add-button\" type=\"submit\" class=\"w-full rounded-md px-4 py-2 text-white transition duration-300\" :class=\"\n            {\n                &#39;bg-blue-500 hover:bg-blue-700&#39;: expenseValue &amp;&amp; !adding,\n                &#39;bg-blue-200&#39;: !expenseValue,\n                &#39;bg-gray-300&#39;: adding\n            }\n        \" :disabled=\"adding || !expenseValue\"><span x-show=\"adding\">Registando...</span> <span x-show=\"!adding\">Registar</span></button> <span class=\"error-message mt-4 block text-rose-500 empty:hidden\"></span></form><div class=\"mt-8\"></div><h2 class=\"mb-4 px-4 text-xl font-bold\">Despesas</h2><ul id=\"expenses-list\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			for _, expense := range expenses {
				templ_7745c5c3_Err = ExpensesListItem(expense).Render(ctx, templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			if len(expenses) == 0 {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<li class=\"px-4\">Nenhuma despesa (ainda)</li>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</ul>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if !templ_7745c5c3_IsBuffer {
				_, templ_7745c5c3_Err = io.Copy(templ_7745c5c3_W, templ_7745c5c3_Buffer)
			}
			return templ_7745c5c3_Err
		})
		templ_7745c5c3_Err = l.Layout().Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
