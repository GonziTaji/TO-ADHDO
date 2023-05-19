import { TaskWithTags } from '@/prismaUtils';

export function TaskList({ tasks }: { tasks: TaskWithTags[] }) {
    return (
        <div className="">
            <h2>Task list</h2>
            <ul className="flex flex-col pt-2">
                {tasks.map((task, i) => (
                    <li key={i} className="border-b border-black">
                        <span className="px-2 py-1">{task.name}</span>

                        <ul className="ms-4 flex">
                            {task.Tags.map((tag) => (
                                <li
                                    key={tag.id}
                                    className="border border-rose-400 rounded bg-amber-300 px-2 m-1"
                                >
                                    {tag.name}
                                </li>
                            ))}

                            {!task.Tags.length && (
                                <li className="border border-rose-300 rounded bg-amber-200 text-gray-600 px-2 m-1">
                                    No tags
                                </li>
                            )}
                        </ul>
                    </li>
                ))}
            </ul>
        </div>
    );
}
