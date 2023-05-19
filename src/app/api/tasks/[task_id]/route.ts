import { NextResponse } from 'next/server';
import { prisma } from '@/prismaUtils';

interface RouteParams {
    params: {
        task_id: string;
    };
}

export async function DELETE(request: Request, { params }: RouteParams) {
    const task_id = parseInt(params.task_id);

    try {
        const tasksResponse = await prisma.tasks.delete({
            where: { id: task_id },
        });

        console.log(tasksResponse);

        return NextResponse.json({ ok: true, data: tasksResponse });
    } catch (error) {
        return new NextResponse(JSON.stringify({ error }), { status: 500 });
    }
}
