export type Applications = {
  items: Array<{
    metadata: { name: string }
  }>;
}

export type Config = {
  allowed: string[];
  refreshInterval: number;
  selectors: {
    environment: string;
  }
}
