import MaterialIcons from '@expo/vector-icons/MaterialIcons'
import { Tabs } from 'expo-router'

export default function TabLayout() {
  return (
    <Tabs screenOptions={{ headerShown: true }}>
      <Tabs.Screen
        name="index"
        options={{
          title: 'Home',
          tabBarIcon: ({ color, size }) => (
            <MaterialIcons name="home" size={size ?? 24} color={color} />
          ),
        }}
      />
      <Tabs.Screen
        name="second"
        options={{
          title: 'Second',
          tabBarIcon: ({ color, size }) => (
            <MaterialIcons name="list" size={size ?? 24} color={color} />
          ),
        }}
      />
    </Tabs>
  )
}
