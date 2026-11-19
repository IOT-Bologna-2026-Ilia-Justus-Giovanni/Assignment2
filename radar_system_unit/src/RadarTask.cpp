#include "RadarTask.h"
#include "Arduino.h"

const unsigned long SEND_INTERVAL_MS = 1000;

void RadarTask::begin() {
  state = State::Idle;
  signalValue = 0;
  lastSentAt = millis();
  Serial.println(signalValue);
}

void RadarTask::update() {
  switch (state) {
    case State::Idle: {
      const unsigned long now = millis();
      if (now - lastSentAt >= SEND_INTERVAL_MS) {
        lastSentAt = now;
        signalValue = signalValue == 0 ? 1 : 0;
        Serial.println(signalValue);
      }
      break;
    }
  }
}
