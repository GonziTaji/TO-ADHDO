'use client';

import { TodoList, TodoListWithTasks } from '@/prismaUtils';
import { Tasks } from '@prisma/client';
import { useState } from 'react';

interface TodoListProps {
    tasks: Tasks[];
    todoList: TodoListWithTasks;
}

export function TodoListView({ tasks, todoList }: TodoListProps) {
    // {
    //     ...todoList,
    //     Tasks: todoList?.Tasks.map(({ Task, ...t }) => ({
    //         ...t,
    //         ...Task,
    //     })),
    // },

    const [selectedTask, setSelectedTask] = useState('');
    const [taskList, setTaskList] = useState<Tasks[]>([]);

    function addTask() {}

    return (
        <div>
            <p>List: {todoList.name}</p>

            <div className="flex gap-4">
                <label htmlFor="task-select">Add a Task: </label>

                <select
                    id="task-select"
                    value={selectedTask}
                    onChange={(ev) => setSelectedTask(ev.currentTarget.value)}
                >
                    {tasks.map((task) => (
                        <option key={task.id}>{task.name}</option>
                    ))}
                </select>
            </div>

            <section>
                <ul>
                    <li></li>
                </ul>
            </section>
        </div>
    );
}
