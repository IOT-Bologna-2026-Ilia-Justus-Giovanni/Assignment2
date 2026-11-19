#ifndef RADAR_TASK_H
#define RADAR_TASK_H

class RadarTask {
public:
  void begin();
  void update();

private:
  enum class State {
    Idle
  };

  State state;
};

#endif
