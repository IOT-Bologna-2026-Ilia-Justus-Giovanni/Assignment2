#include "RadarTask.h"

void RadarTask::begin() {
  state = State::Idle;
}

void RadarTask::update() {
  switch (state) {
    case State::Idle:
      break;
  }
}
