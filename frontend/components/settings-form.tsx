"use client"

import type { ImageSettings } from "@/app/page"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Label } from "@/components/ui/label"
import { Input } from "@/components/ui/input"
import { Slider } from "@/components/ui/slider"

interface SettingsFormProps {
  settings: ImageSettings
  onSettingsChange: (settings: ImageSettings) => void
}

export function SettingsForm({ settings, onSettingsChange }: SettingsFormProps) {
  const handleWidthChange = (value: string) => {
    const width = Number.parseInt(value)
    if (!isNaN(width) && width >= 0) {
      onSettingsChange({ ...settings, maxWidth: width })
    }
  }

  const handleHeightChange = (value: string) => {
    const height = Number.parseInt(value)
    if (!isNaN(height) && height >= 0) {
      onSettingsChange({ ...settings, maxHeight: height })
    }
  }

  const handleQualityChange = (value: number[]) => {
    onSettingsChange({ ...settings, quality: value[0] })
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Image Settings</CardTitle>
      </CardHeader>
      <CardContent className="space-y-6">
        <div className="space-y-2">
          <Label htmlFor="max-width">Max Width (px)</Label>
          <Input
            id="max-width"
            type="number"
            min="0"
            value={settings.maxWidth}
            onChange={(e) => handleWidthChange(e.target.value)}
          />
        </div>

        <div className="space-y-2">
          <Label htmlFor="max-height">Max Height (px)</Label>
          <Input
            id="max-height"
            type="number"
            min="0"
            value={settings.maxHeight}
            onChange={(e) => handleHeightChange(e.target.value)}
          />
        </div>

        <div className="space-y-4">
          <div className="flex justify-between">
            <Label htmlFor="quality">Quality</Label>
            <span className="text-sm font-medium">{settings.quality}%</span>
          </div>
          <Slider
            id="quality"
            min={1}
            max={100}
            step={1}
            value={[settings.quality]}
            onValueChange={handleQualityChange}
          />
          <div className="flex justify-between text-xs text-muted-foreground">
            <span>Lower quality, smaller file</span>
            <span>Higher quality, larger file</span>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}
