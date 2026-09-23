import { afterEach } from "vitest";
import { cleanup } from "@testing-library/react";
import "@testing-library/jest-dom/vitest";

// Necessário explicitamente com `isolate: false` (ver vite.config.js): sem
// isso, o DOM de um teste pode continuar montado no próximo teste do mesmo
// arquivo, dando falso positivo de "elemento duplicado encontrado".
afterEach(() => {
  cleanup();
});
