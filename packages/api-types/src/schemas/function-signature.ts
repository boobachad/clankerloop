import { z } from "@hono/zod-openapi";

// TypeDef represents any type in the schema
// We need to define the type first, then use z.lazy for the recursive schema
export type TypeDef =
  | { kind: "primitive"; type: "int" | "float" | "string" | "boolean" | "null" }
  | { kind: "array"; items: TypeDef }
  | { kind: "object"; properties: Record<string, TypeDef> }
  | { kind: "map"; keyType: TypeDef; valueType: TypeDef }
  | { kind: "tuple"; items: TypeDef[] }
  | { kind: "union"; types: TypeDef[] }
  | { kind: "reference"; name: string };

// Recursive type definition for function parameter/return types
export const TypeDefSchema: z.ZodType<TypeDef> = z.lazy(() =>
  z.discriminatedUnion("kind", [
    z.object({
      kind: z.literal("primitive"),
      type: z.enum(["int", "float", "string", "boolean", "null"]),
    }),
    z.object({
      kind: z.literal("array"),
      items: TypeDefSchema,
    }),
    z.object({
      kind: z.literal("object"),
      properties: z.record(TypeDefSchema),
    }),
    z.object({
      kind: z.literal("map"),
      keyType: TypeDefSchema,
      valueType: TypeDefSchema,
    }),
    z.object({
      kind: z.literal("tuple"),
      items: z.array(TypeDefSchema),
    }),
    z.object({
      kind: z.literal("union"),
      types: z.array(TypeDefSchema),
    }),
    // Reference to a named type (for recursive types like TreeNode)
    z.object({
      kind: z.literal("reference"),
      name: z.string(),
    }),
  ]),
);

// Named type definition (for recursive/reusable types)
export const NamedTypeSchema = z.object({
  name: z.string(),
  definition: TypeDefSchema,
  description: z.string().optional(),
});

export const FunctionParameterSchema = z.object({
  name: z.string(),
  type: TypeDefSchema,
  optional: z.boolean().optional(),
  description: z.string().optional(),
});

// Complete function signature schema
export const FunctionSignatureSchemaSchema = z
  .object({
    version: z.literal(1),
    functionName: z.string().default("runSolution"),
    parameters: z.array(FunctionParameterSchema),
    returnType: TypeDefSchema,
    // Named types that can be referenced (for recursive types like TreeNode)
    namedTypes: z.array(NamedTypeSchema).optional(),
  })
  .openapi("FunctionSignatureSchema");

// Inferred types
export type NamedType = z.infer<typeof NamedTypeSchema>;
export type FunctionParameter = z.infer<typeof FunctionParameterSchema>;
export type FunctionSignatureSchema = z.infer<
  typeof FunctionSignatureSchemaSchema
>;
