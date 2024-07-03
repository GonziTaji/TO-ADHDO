import { getTasksWithTagsOfUser, prisma } from '@/prismaUtils';
import BackButton from './BackButton';
import { TodoListView } from './TodoListView';

const user_id = 1;

async function getTodoList(todo_list_id: string) {
    try {
        const todoList = await prisma.todoLists.findFirst({
            where: { id: parseInt(todo_list_id) },
            include: { Tasks: { include: { Task: true } } },
        });

        const tasks = await getTasksWithTagsOfUser(user_id);

        return { todoList, tasks };
    } catch (error) {
        return { error };
    }
}

interface PageProps {
    params: {
        todo_list_id: string;
    };
}

export default async function Page({ params: { todo_list_id } }: PageProps) {
    const { todoList, tasks, error } = await getTodoList(todo_list_id);

    if (!todoList || !todoList.id) {
        return (
            <div>
                <BackButton />
                <p>404 - todo list not found</p>
            </div>
        );
    }

    if (error) {
        console.error(error);
    }

    return (
        <div>
            <BackButton />

            <TodoListView tasks={tasks} todoList={todoList} />
        </div>
    );
}
