import React from "react";

import { Footer } from "rsuite";

export default function AppFooter() {
  return (
    <Footer
      style={{ textAlign: "center", padding: 15, background: "#00000006" }}
    >
      Copyright © 2020 - {new Date().getFullYear()} <b>ZMIDE Studio.</b>
    </Footer>
  );
}
