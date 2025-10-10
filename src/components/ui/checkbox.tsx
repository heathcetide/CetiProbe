import { twMerge } from "tailwind-merge";
import { Check } from "lucide-react";
import { Checkbox } from "@ark-ui/react/checkbox";

function CheckboxRoot({
  children,
  className,
  ...rest
}: React.ComponentProps<typeof Checkbox.Root>) {
  return (
    <Checkbox.Root
      className={twMerge([
        "flex items-center space-x-2",
        className,
      ])}
      {...rest}
    >
      {children}
    </Checkbox.Root>
  );
}

function CheckboxLabel({
  children,
  className,
  ...rest
}: React.ComponentProps<typeof Checkbox.Label>) {
  return (
    <Checkbox.Label 
      className={twMerge([
        "text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70",
        className
      ])} 
      {...rest}
    >
      {children}
    </Checkbox.Label>
  );
}

function CheckboxControl({
  className,
  ...rest
}: React.ComponentProps<typeof Checkbox.Control>) {
  return (
    <Checkbox.Control
      className={twMerge([
        "peer h-4 w-4 shrink-0 rounded-sm border border-primary ring-offset-background",
        "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2",
        "disabled:cursor-not-allowed disabled:opacity-50",
        "data-[state=checked]:bg-primary data-[state=checked]:text-primary-foreground",
        className,
      ])}
      {...rest}
    >
      <Check className="h-4 w-4" />
    </Checkbox.Control>
  );
}

function CheckboxHiddenInput() {
  return <Checkbox.HiddenInput />;
}

export { CheckboxRoot, CheckboxLabel, CheckboxControl, CheckboxHiddenInput };
