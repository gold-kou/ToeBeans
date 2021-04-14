import React from "react";
import { Button } from "react-bootstrap";
import RefreshIcon from "@material-ui/icons/Refresh";
import "./LoaderButton.css";

export default function LoaderButton({
  isLoading,
  className = "",
  disabled = false,
  ...props
}) {
  return (
    <Button
      className={`LoaderButton ${className}`}
      disabled={disabled || isLoading}
      {...props}
    >
      {isLoading && <RefreshIcon></RefreshIcon>}
      {props.children}
    </Button>
  );
}
