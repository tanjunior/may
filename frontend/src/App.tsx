import reactLogo from "./assets/react.svg";
import viteLogo from "/vite.svg";
import "./App.css";
import { Suspense } from "react";
import { ModeToggle } from "./components/mode-toggle";
import LatestProduct from "./components/LatestProduct";
import CRUDExamples from "./components/CRUDExamples";
import { ThemeProvider } from "./components/theme-provider";

function App() {
  return (
    <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
      <div>
        <header>
          <div>
            <a href="/" className="flex items-center gap-3">
              <img src={viteLogo} className="logo" alt="Vite" />
              <div>
                <div>Vite + React</div>
                <div className="text-sm text-slate-500">
                  CRUD demo · shadcn style
                </div>
              </div>

            </a>
          </div>
              <ModeToggle />
        </header>

        <main>
          <section>
            <div>
              <Suspense>
                <LatestProduct />
              </Suspense>
            </div>
          </section>

          <section>
            <div>
              <CRUDExamples />
            </div>
          </section>
        </main>
      </div>
    </ThemeProvider>
  );
}

export default App;
