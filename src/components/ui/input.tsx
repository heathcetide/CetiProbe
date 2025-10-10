import { twMerge } from "tailwind-merge";
import { forwardRef } from "react";

const Input = forwardRef<HTMLInputElement, React.ComponentProps<"input">>(
  ({ className, type, ...props }, ref) => {
    return (
      <div className="relative group">
        <input
          type={type}
          className={twMerge(
            "flex h-10 w-full rounded-lg border border-white/20 bg-white/5 backdrop-blur-sm px-4 py-2 text-sm text-white placeholder:text-white/60",
            "transition-all duration-300 ease-out",
            "focus:outline-none focus:border-white/40 focus:bg-white/10 focus:shadow-lg focus:shadow-white/10",
            "hover:border-white/30 hover:bg-white/8",
            "disabled:cursor-not-allowed disabled:opacity-50",
            "group-hover:shadow-md group-hover:shadow-white/5",
            className
          )}
          ref={ref}
          {...props}
        />
        <div className="absolute inset-0 rounded-lg bg-gradient-to-r from-white/10 to-transparent opacity-0 group-focus-within:opacity-100 transition-opacity duration-300 pointer-events-none" />
      </div>
    );
  }
);
Input.displayName = "Input";

export { Input };
