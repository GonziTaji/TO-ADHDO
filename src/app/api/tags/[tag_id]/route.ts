import { NextResponse } from 'next/server';
import { prisma } from '@/prismaUtils';
import { validateParams } from '@/utils';

interface RouteParams {
    params: {
        tag_id: string;
    };
}

export async function DELETE(request: Request, { params }: RouteParams) {
    const tag_id = parseInt(params.tag_id);

    try {
        const tagResponse = await prisma.tags.delete({
            where: { id: tag_id },
        });

        console.log(tagResponse);

        return NextResponse.json({ ok: true, data: tagResponse });
    } catch (error: any) {
        return new NextResponse(JSON.stringify({ error: error.toString() }), {
            status: 500,
        });
    }
}

export async function PUT(request: Request, { params }: RouteParams) {
    const tag_id = parseInt(params.tag_id);

    const body = await request.json();

    const error = validateParams(body, ['name']);

    if (error) {
        return new NextResponse(JSON.stringify({ message: error }), {
            status: 400,
        });
    }

    try {
        const tagResponse = await prisma.tags.update({
            where: { id: tag_id },
            data: {
                name: body.name,
            },
        });

        console.log(tagResponse);

        return NextResponse.json({ ok: true, data: tagResponse });
    } catch (error: any) {
        return new NextResponse(JSON.stringify({ error: error.toString() }), {
            status: 500,
        });
    }
}
