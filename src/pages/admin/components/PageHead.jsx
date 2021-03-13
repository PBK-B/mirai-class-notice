import React from "react";
import PropTypes from "prop-types";

import { Divider, Row, Button, Icon } from "rsuite";

export default function PageHead(props) {
  const { title, onLeftClick } = props;
  return (
    <div>
      <Row>
        <Button appearance="link" onClick={onLeftClick}>
          <Icon icon="angle-left" style={{ fontSize: 22 }} />
        </Button>
        <Divider vertical />
        {title && <h6 style={{ display: "inline-flex" }}>{title}</h6>}
      </Row>
      <Divider />
    </div>
  );
}

PageHead.prototype = {
  title: PropTypes.string,
  onLeftClick: PropTypes.func,
};
