export function whichSidesInt(...i) {
  let top,
    right,
    bottom,
    left,
    ok = false;

  switch (i.length) {
    case 1:
      top = i[0];
      bottom = i[0];
      left = i[0];
      right = i[0];
      ok = true;
      break;
    case 2:
      top = i[0];
      bottom = i[0];
      left = i[1];
      right = i[1];
      ok = true;
      break;
    case 3:
      top = i[0];
      left = i[1];
      right = i[1];
      bottom = i[2];
      ok = true;
      break;
    case 4:
      top = i[0];
      right = i[1];
      bottom = i[2];
      left = i[3];
      ok = true;
      break;
  }

  return [top, right, bottom, left, ok];
}

export function whichSidesBool(...i) {
  let top,
    right,
    bottom,
    left,
    ok = false;

  switch (i.length) {
    case 1:
      top = i[0];
      bottom = i[0];
      left = i[0];
      right = i[0];
      ok = true;
      break;
    case 2:
      top = i[0];
      bottom = i[0];
      left = i[1];
      right = i[1];
      ok = true;
      break;
    case 3:
      top = i[0];
      left = i[1];
      right = i[1];
      bottom = i[2];
      ok = true;
      break;
    case 4:
      top = i[0];
      right = i[1];
      bottom = i[2];
      left = i[3];
      ok = true;
      break;
  }

  return [top, right, bottom, left, ok];
}

export function whichSidesColor(...colors) {
  let top,
    right,
    bottom,
    left,
    ok = false;

  switch (colors.length) {
    case 1:
      top = colors[0];
      bottom = colors[0];
      left = colors[0];
      right = colors[0];
      ok = true;
      break;
    case 2:
      top = colors[0];
      bottom = colors[0];
      left = colors[1];
      right = colors[1];
      ok = true;
      break;
    case 3:
      top = colors[0];
      left = colors[1];
      right = colors[1];
      bottom = colors[2];
      ok = true;
      break;
    case 4:
      top = colors[0];
      right = colors[1];
      bottom = colors[2];
      left = colors[3];
      ok = true;
      break;
  }

  return [top, right, bottom, left, ok];
}
