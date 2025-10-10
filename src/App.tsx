import { createContext, useState } from "react";
import { twMerge } from "tailwind-merge";
import { Outlet } from "react-router";
import { UnifiedHeader } from "@/components/UnifiedHeader";

export const MobileMenuContext = createContext<{
    showMenu: boolean;
    setShowMenu: React.Dispatch<React.SetStateAction<boolean>>;
}>({
    showMenu: true,
    setShowMenu: () => {},
});

function App() {
    const [showMenu, setShowMenu] = useState(false);

    return (
        <MobileMenuContext.Provider value={{ showMenu, setShowMenu }}>
            <div
                className={twMerge([
                    "min-h-screen",
                    "before:fixed before:inset-0 before:bg-noise before:z-[-1]",
                    "after:bg-temper after:opacity-15 after:bg-contain after:fixed after:inset-0 after:blur-xl after:z-[-1]",
                ])}
            >
                {/* 统一Header */}
                <UnifiedHeader 
                    showMenu={showMenu} 
                    setShowMenu={setShowMenu}
                />
                
                {/* 主要内容区域 */}
                <div className="pt-16">
                    <div className="container mx-auto px-4 sm:px-6">
                        <Outlet />
                        <div className="text-center text-foreground/60 mt-24 mb-20">
                            Powered by synthetic caffeine ∙ Deployed by{" "}
                            <a
                                href="https://github.com/heathcetide"
                                target="_blank"
                                className="font-medium text-foreground/70 hover:text-foreground transition-colors"
                            >
                                Left4code
                            </a>{" "}
                            ∙ Signal traceable on{" "}
                            <a
                                href="https://github.com/heathcetide"
                                target="_blank"
                                className="font-medium text-foreground/70 hover:text-foreground transition-colors"
                            >
                                GitHub
                            </a>
                            .
                        </div>
                    </div>
                </div>
            </div>
        </MobileMenuContext.Provider>
    );
}

export default App;
