export type RuneFunc = (input: string) => Promise<string>;
export type Middleware = (name: string, next: RuneFunc) => RuneFunc;

export class Dispatcher {
  private readonly runes: Map<string, RuneFunc>;
  private readonly middleware: Middleware[];

  constructor(runes: Map<string, RuneFunc>, middleware: Middleware[]) {
    this.runes = new Map(runes);
    this.middleware = [...middleware];
  }

  async call(name: string, input: string): Promise<string> {
    let fn = this.runes.get(name);
    if (!fn) {
      throw new Error(`callable "${name}" not registered`);
    }
    for (let i = this.middleware.length - 1; i >= 0; i--) {
      fn = this.middleware[i](name, fn);
    }
    return fn(input);
  }
}
