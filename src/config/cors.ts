import cors from "@elysiajs/cors";

export default function corsConfig() {
  return cors({
    credentials: true,
    maxAge: 3600 * 24,
    origin: true,
    methods: ["GET", "PATCH", "DELETE", "POST"],
  });
}
