import { twMerge } from "tailwind-merge";
import { Button } from "@/components/ui/button";
import { X } from "lucide-react";
import { Dialog } from "@ark-ui/react/dialog";
import { Portal } from "@ark-ui/react/portal";

function DialogRoot({
  children,
  ...rest
}: React.ComponentProps<typeof Dialog.Root>) {
  return <Dialog.Root {...rest}>{children}</Dialog.Root>;
}

function DialogTrigger({
  children,
  className,
  ...rest
}: React.ComponentProps<typeof Dialog.Trigger>) {
  return (
    <Dialog.Trigger asChild {...rest}>
      <Button className={className}>{children}</Button>
    </Dialog.Trigger>
  );
}

function DialogBackdrop({
  className,
  ...rest
}: React.ComponentProps<typeof Dialog.Backdrop>) {
  return (
    <Dialog.Backdrop
      className={twMerge([
        "fixed inset-0 bg-background/80 backdrop-blur-sm z-50",
        "data-[state=open]:animate-in data-[state=open]:fade-in-0",
        "data-[state=closed]:animate-out data-[state=closed]:fade-out-0",
        className,
      ])}
      {...rest}
    />
  );
}

function DialogPositioner({
  children,
  className,
  ...rest
}: React.ComponentProps<typeof Dialog.Positioner>) {
  return (
    <Dialog.Positioner className={className} {...rest}>
      {children}
    </Dialog.Positioner>
  );
}

function DialogContent({
  children,
  className,
  ...rest
}: React.ComponentProps<typeof Dialog.Content>) {
  return (
    <Dialog.Content
      className={twMerge([
        "fixed left-[50%] top-[50%] z-50 grid w-full max-w-lg translate-x-[-50%] translate-y-[-50%] gap-4 border bg-background p-6 shadow-lg duration-200",
        "data-[state=open]:animate-in data-[state=open]:fade-in-0 data-[state=open]:zoom-in-95 data-[state=open]:slide-in-from-left-1/2 data-[state=open]:slide-in-from-top-[48%]",
        "data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=closed]:zoom-out-95 data-[state=closed]:slide-out-to-left-1/2 data-[state=closed]:slide-out-to-top-[48%]",
        "sm:rounded-lg",
        className,
      ])}
      {...rest}
    >
      {children}
    </Dialog.Content>
  );
}

function DialogTitle({
  children,
  className,
  ...rest
}: React.ComponentProps<typeof Dialog.Title>) {
  return (
    <Dialog.Title
      className={twMerge([
        "text-lg font-semibold leading-none tracking-tight",
        className,
      ])}
      {...rest}
    >
      {children}
    </Dialog.Title>
  );
}

function DialogDescription({
  children,
  className,
  ...rest
}: React.ComponentProps<typeof Dialog.Description>) {
  return (
    <Dialog.Description
      className={twMerge(["text-sm text-muted-foreground", className])}
      {...rest}
    >
      {children}
    </Dialog.Description>
  );
}

function DialogCloseTrigger({
  children,
  className,
  asChild,
  ...rest
}: React.ComponentProps<typeof Dialog.CloseTrigger>) {
  return (
    <Dialog.CloseTrigger asChild {...rest}>
      {!asChild ? (
        <Button
          className={twMerge([
            "absolute right-4 top-4 rounded-sm opacity-70 ring-offset-background transition-opacity hover:opacity-100 focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:pointer-events-none",
            className,
          ])}
          {...rest}
        >
          <X className="h-4 w-4" />
        </Button>
      ) : (
        children
      )}
    </Dialog.CloseTrigger>
  );
}

export {
  DialogRoot,
  DialogTrigger,
  DialogBackdrop,
  DialogPositioner,
  DialogContent,
  DialogTitle,
  DialogDescription,
  DialogCloseTrigger,
  Portal,
};
