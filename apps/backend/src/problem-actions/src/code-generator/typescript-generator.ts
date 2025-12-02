import type { FunctionSignatureSchema, TypeDef } from "@repo/api-types";
import type { CodeGenerator } from "./types";

export class TypeScriptGenerator implements CodeGenerator {
  typeToString(typeDef: TypeDef): string {
    switch (typeDef.kind) {
      case "primitive":
        if (typeDef.type === "int" || typeDef.type === "float") {
          return "number";
        }
        return typeDef.type;
      case "array":
        return `${this.typeToString(typeDef.items)}[]`;
      case "object": {
        const props = Object.entries(typeDef.properties)
          .map(([k, v]) => `${k}: ${this.typeToString(v)}`)
          .join("; ");
        return `{ ${props} }`;
      }
      case "map":
        return `Record<${this.typeToString(typeDef.keyType)}, ${this.typeToString(typeDef.valueType)}>`;
      case "tuple":
        return `[${typeDef.items.map((i) => this.typeToString(i)).join(", ")}]`;
      case "union":
        return typeDef.types.map((t) => this.typeToString(t)).join(" | ");
      case "reference":
        return typeDef.name;
    }
  }

  generateTypeDefinitions(schema: FunctionSignatureSchema): string {
    if (!schema.namedTypes?.length) return "";

    return schema.namedTypes
      .map((nt) => {
        if (nt.definition.kind === "object") {
          const props = Object.entries(nt.definition.properties)
            .map(([k, v]) => `  ${k}: ${this.typeToString(v)};`)
            .join("\n");
          return `interface ${nt.name} {\n${props}\n}`;
        }
        return `type ${nt.name} = ${this.typeToString(nt.definition)};`;
      })
      .join("\n\n");
  }

  generateScaffold(schema: FunctionSignatureSchema): string {
    const params = schema.parameters
      .map((p) => {
        const optional = p.optional ? "?" : "";
        return `${p.name}${optional}: ${this.typeToString(p.type)}`;
      })
      .join(", ");
    const returnType = this.typeToString(schema.returnType);
    return `function runSolution(${params}): ${returnType} {\n  // TODO: implement your solution here\n  throw new Error("Not implemented");\n}`;
  }

  generateStarterCode(schema: FunctionSignatureSchema): string {
    const typeDefs = this.generateTypeDefinitions(schema);
    const scaffold = this.generateScaffold(schema);

    return typeDefs ? `${typeDefs}\n\n${scaffold}` : scaffold;
  }
}
