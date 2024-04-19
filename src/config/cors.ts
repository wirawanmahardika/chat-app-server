import cors from "@elysiajs/cors";

export default function corsConfig() {
  return cors({
    credentials: true,
    maxAge: 3600 * 24,
    origin: ["localhost:5173", "127.0.0.1:5500", "localhost:5500"],
    methods: ["GET", "PATCH", "DELETE", "POST"],
    allowedHeaders: ["Content-Type"],
    preflight: true,
  });
}
