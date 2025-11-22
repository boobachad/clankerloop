const BRAINTRUST_API_KEY = process.env.BRAINTRUST_API_KEY;

if (!BRAINTRUST_API_KEY) {
  throw new Error("BRAINTRUST_API_KEY is not set");
}

export const env = {
  BRAINTRUST_API_KEY,
};
