import type { ActionFunctionArgs, MetaFunction } from "@remix-run/node";
import {
  Form,
  json,
  useActionData,
  useLoaderData,
  useNavigation,
} from "@remix-run/react";
import { db } from "~/modules/db";
import { DateTime } from "luxon";
import React from "react";

export const meta: MetaFunction = () => {
  return [
    { title: "Mercado" },
    { name: "description", content: "Welcome to Remix!" },
  ];
};

const MONTH_START = DateTime.local({
  zone: "Europe/Lisbon",
  locale: "pt",
}).startOf("month");
const MONTH_END = MONTH_START.endOf("month");
const GROCERY_STORES = [
  "Continente",
  "Froiz",
  "Lidl",
  "Mercadona",
  "Minipreço",
  "Pingo Doce",
];

export async function action({ request }: ActionFunctionArgs) {
  const formData = await request.formData();
  const amount = Number(formData.get("amount")) * 100;

  if (amount <= 0) {
    return json({ error: "O valor deve ser maior do que zero." });
  }

  const groceryStore = formData.get("groceryStore");
  const expense = await db.expense.create({
    data: {
      amount,
      groceryStore: String(groceryStore),
    },
  });

  return json({ expense });
}

export async function loader() {
  const expenses = await db.expense.findMany({
    where: {
      created_at: {
        gte: MONTH_START.toJSDate(),
        lte: MONTH_END.toJSDate(),
      },
    },
    orderBy: {
      created_at: "desc",
    },
  });

  const total = expenses.reduce((acc, expense) => acc + expense.amount, 0);
  const average = total / expenses.length;
  const currentMonth = MONTH_START.toFormat("LLLL");
  const year = MONTH_START.toFormat("yyyy");
  const daysInMonth = MONTH_START.daysInMonth;
  const today = DateTime.local({ zone: "Europe/Lisbon", locale: "pt" });
  const monthPercentage = ((today.day - 1) / daysInMonth) * 100;

  return json({
    expenses,
    total,
    average,
    currentMonth,
    monthPercentage,
    year,
  });
}

function formatEuros(amount: number) {
  return `${(amount / 100).toFixed(2)} €`;
}

function DecimalInput(props: React.InputHTMLAttributes<HTMLInputElement>) {
  const [value, set_value] = React.useState<string>("0.00");

  const formatInputValue = (inputValue: string): string => {
    const numericValue = inputValue.replace(/\D/g, "");
    const decimalValue = Number(numericValue) / 100;
    return decimalValue.toFixed(2);
  };

  const handleChange = (event: React.ChangeEvent<HTMLInputElement>): void => {
    const formattedValue = formatInputValue(event.target.value);
    set_value(formattedValue);
  };

  return (
    <input
      {...props}
      type="text"
      inputMode="decimal"
      value={value}
      onChange={handleChange}
    />
  );
}

export default function Index() {
  const { expenses, total, average, currentMonth, year, monthPercentage } =
    useLoaderData<typeof loader>();
  const navigation = useNavigation();
  useLoaderData<typeof loader>();
  const formRef = React.useRef<HTMLFormElement>(null);
  const actionData = useActionData<typeof action>();

  React.useEffect(() => {
    if (navigation.state !== "idle") {
      return;
    }

    const input = document.querySelector<HTMLInputElement>(
      "input[name='amount']"
    );

    if (input == null) {
      return;
    }

    input.value = "0.00";
  }, [navigation.state]);

  return (
    <main className="font-sans p-4">
      <h1 className="text-2xl">
        {currentMonth} {year} | {Math.round(monthPercentage)}%
      </h1>
      {`Total: ${formatEuros(total)}`}
      <br />
      {`Média: ${formatEuros(average)}`}

      <Form method="post" ref={formRef} className="my-4 space-y-2">
        <label htmlFor="amount">Valor</label>
        <DecimalInput
          name="amount"
          className="w-full rounded border border-gray-300 p-2"
        />

        {actionData != null && "error" in actionData ? (
          <p className="text-red-500">{actionData.error}</p>
        ) : null}

        <label htmlFor="groceryStore">Mercado</label>
        <select
          name="groceryStore"
          className="w-full rounded border border-gray-300 p-2"
        >
          {GROCERY_STORES.map((groceryStore) => (
            <option key={groceryStore} value={groceryStore}>
              {groceryStore}
            </option>
          ))}
        </select>

        <div className="mt-4" />

        <button
          type="submit"
          className="rounded w-full bg-blue-500 py-2 px-4 text-white"
        >
          Registar
        </button>
      </Form>

      <div>
        {expenses.map((expense) => {
          const date = DateTime.fromISO(expense.created_at, {
            zone: "Europe/Lisbon",
            locale: "pt",
          });
          const formattedDate = date.toFormat("EEE dd LLL T");

          return (
            <div key={expense.id}>
              <div>
                {formatEuros(expense.amount)} - {formattedDate}
              </div>
              <div>{expense.groceryStore}</div>
            </div>
          );
        })}
      </div>
    </main>
  );
}
