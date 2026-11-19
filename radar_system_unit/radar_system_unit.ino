#include "RadarTask.h"

RadarTask radarTask;

void setup() {
  Serial.begin(9600);
  radarTask.begin();
}

void loop() {
  radarTask.update();
}
