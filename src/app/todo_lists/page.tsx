import { TodoList, getTodoLists } from '@/prismaUtils';
import NewTodoListForm from './NewTodoListForm';
import TodoListItemList from './TodoListItemList';

async function _getTodoLists() {
    try {
        const todoLists = await getTodoLists(1);

        return { todoLists };
    } catch (error) {
        return { error, todoLists: [] as TodoList[] };
    }
}

export default async function Page() {
    const { todoLists, error } = await _getTodoLists();

    if (error) {
        console.error(error);
    }

    return (
        <div>
            <h1>To-Do Lists</h1>
            <section>
                <h2>Create a new List</h2>
                <NewTodoListForm />
            </section>

            <section>
                <h2>Select a List</h2>

                <ul>
                    {todoLists.map((todo) => (
                        <li key={todo.id}>
                            <TodoListItemList todoList={todo} />
                        </li>
                    ))}
                </ul>
            </section>
        </div>
    );
}
