import { db } from '@vercel/postgres';
import { NextApiRequest, NextApiResponse } from 'next';
 
export default async function tasks(
  request: NextApiRequest,
  response: NextApiResponse,
) {
  const client = await db.connect();
 
  try {
    const tasksResponse = await client.sql`SELECT * FROM tasks;`;
    return response.status(200).json({ tasks: tasksResponse.fields });
  } catch (error) {
    return response.status(500).json({ error });
  }
}