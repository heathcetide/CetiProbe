import { twMerge } from "tailwind-merge";
import { Switch } from "@ark-ui/react/switch";

function SwitchRoot({
  className,
  children,
  ...rest
}: React.ComponentProps<typeof Switch.Root>) {
  return (
    <Switch.Root
      className={twMerge(["flex items-center gap-3", className])}
      {...rest}
    >
      {children}
    </Switch.Root>
  );
}

function SwitchHiddenInput() {
  return <Switch.HiddenInput />;
}

function SwitchControl({
  className,
  children,
  ...rest
}: React.ComponentProps<typeof Switch.Control>) {
  return (
    <Switch.Control
      className={twMerge([
        "peer inline-flex h-6 w-11 shrink-0 cursor-pointer items-center rounded-full border-2 border-transparent transition-colors",
        "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background",
        "disabled:cursor-not-allowed disabled:opacity-50",
        "data-[state=checked]:bg-primary data-[state=unchecked]:bg-input",
        className,
      ])}
      {...rest}
    >
      {children}
    </Switch.Control>
  );
}

function SwitchThumb({
  className,
  ...rest
}: React.ComponentProps<typeof Switch.Thumb>) {
  return (
    <Switch.Thumb
      className={twMerge([
        "pointer-events-none block h-5 w-5 rounded-full bg-background shadow-lg ring-0 transition-transform",
        "data-[state=checked]:translate-x-5 data-[state=unchecked]:translate-x-0",
        className,
      ])}
      {...rest}
    />
  );
}

function SwitchLabel({
  className,
  children,
  ...rest
}: React.ComponentProps<typeof Switch.Label>) {
  return (
    <Switch.Label 
      className={twMerge([
        "text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70",
        className
      ])} 
      {...rest}
    >
      {children}
    </Switch.Label>
  );
}

export {
  SwitchRoot,
  SwitchHiddenInput,
  SwitchControl,
  SwitchThumb,
  SwitchLabel,
};
