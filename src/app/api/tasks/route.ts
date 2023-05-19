import { NextResponse } from 'next/server';
import { PrismaClient } from '@prisma/client';

export async function GET(request: Request) {
    return new NextResponse(JSON.stringify({ hola: 'mundo' }));
}

export async function POST(request: Request) {
    const body = await request.json();

    console.log(body);

    if (!body.name) {
        return new NextResponse(
            JSON.stringify({ message: 'task must have a name' }),
            { status: 400 }
        );
    }

    const prisma = new PrismaClient();
    const tagsIds: number[] = Array.isArray(body.tags) ? body.tags : [];

    try {
        const tasksResponse = await prisma.tasks.create({
            data: {
                name: body.name,
                user_id: body.user_id,
                Tags: {
                    connect: tagsIds.map((id) => ({ id })),
                },
            },
        });

        console.log(tasksResponse);

        return new NextResponse(JSON.stringify({ ok: true }));
    } catch (error) {
        return new NextResponse(JSON.stringify({ error }), { status: 500 });
    }
}
