import jwt from "@elysiajs/jwt";
import { t } from "elysia";

export default function jwtConf() {
  return jwt({
    secret: process.env.JWT_SECRET || "aosiudoiweu83579384",
    exp: "1d",
    name: "jwt",
    schema: t.Object({ username: t.String(), id: t.String() }),
  });
}
