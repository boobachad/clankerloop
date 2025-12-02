import type { FunctionSignatureSchema, TypeDef } from "@repo/api-types";
import type { CodeGenerator } from "./types";

export class PythonGenerator implements CodeGenerator {
  private imports = new Set<string>();

  typeToString(typeDef: TypeDef): string {
    switch (typeDef.kind) {
      case "primitive": {
        const map: Record<string, string> = {
          int: "int",
          float: "float",
          string: "str",
          boolean: "bool",
          null: "None",
        };
        return map[typeDef.type];
      }
      case "array":
        return `list[${this.typeToString(typeDef.items)}]`;
      case "object":
        // Inline objects become dict (could use TypedDict for more precision)
        return "dict";
      case "map":
        return `dict[${this.typeToString(typeDef.keyType)}, ${this.typeToString(typeDef.valueType)}]`;
      case "tuple":
        return `tuple[${typeDef.items.map((i) => this.typeToString(i)).join(", ")}]`;
      case "union": {
        // Simplify T | None to Optional[T]
        const nonNullTypes = typeDef.types.filter(
          (t) => !(t.kind === "primitive" && t.type === "null"),
        );
        const hasNull = typeDef.types.some(
          (t) => t.kind === "primitive" && t.type === "null",
        );

        if (hasNull && nonNullTypes.length === 1) {
          this.imports.add("Optional");
          return `Optional[${this.typeToString(nonNullTypes[0])}]`;
        }

        this.imports.add("Union");
        return `Union[${typeDef.types.map((t) => this.typeToString(t)).join(", ")}]`;
      }
      case "reference":
        // Forward reference as string for recursive types
        return `"${typeDef.name}"`;
    }
  }

  generateTypeDefinitions(schema: FunctionSignatureSchema): string {
    if (!schema.namedTypes?.length) return "";

    return schema.namedTypes
      .map((nt) => {
        if (nt.definition.kind === "object") {
          const props = Object.entries(nt.definition.properties)
            .map(([k, v]) => `    ${k}: ${this.typeToString(v)}`)
            .join("\n");
          return `class ${nt.name}:\n${props}`;
        }
        return `${nt.name} = ${this.typeToString(nt.definition)}`;
      })
      .join("\n\n");
  }

  generateScaffold(schema: FunctionSignatureSchema): string {
    const params = schema.parameters
      .map((p) => {
        const typeStr = this.typeToString(p.type);
        if (p.optional) {
          this.imports.add("Optional");
          return `${p.name}: Optional[${typeStr}] = None`;
        }
        return `${p.name}: ${typeStr}`;
      })
      .join(", ");
    const returnType = this.typeToString(schema.returnType);
    return `def run_solution(${params}) -> ${returnType}:\n    # TODO: implement your solution here\n    raise NotImplementedError()`;
  }

  generateStarterCode(schema: FunctionSignatureSchema): string {
    // Clear imports before generating to get fresh set
    this.imports.clear();

    // Generate type definitions first (this populates imports)
    const typeDefs = this.generateTypeDefinitions(schema);

    // Generate scaffold (this also populates imports)
    const scaffold = this.generateScaffold(schema);

    // Build import line
    const importLine =
      this.imports.size > 0
        ? `from typing import ${Array.from(this.imports).sort().join(", ")}\n\n`
        : "";

    // Combine all parts
    const parts = [importLine];
    if (typeDefs) {
      parts.push(typeDefs + "\n\n");
    }
    parts.push(scaffold);

    return parts.join("");
  }
}
